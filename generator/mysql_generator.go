package generator

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/whencome/gomodel"
    "github.com/whencome/gotil/fileutil"
    "gomodelcreator/config"
    "gomodelcreator/filemaker"
    "strings"
)

//-------------------- MODEL DEFINITION ---------------------//
var msMgr = gomodel.NewModelManager(&MySQLSchema{})

// define a empty model
type MySQLSchema struct{}

// GetDatabase 获取数据库名称（返回配置中的名称，不要使用实际数据库名称，因为实际数据库名称在不同环境可能不一样）
func (s *MySQLSchema) GetDatabase() string {
    return "mysqldb"
}

// GetTableName 获取数据库数据存放的数据表名称
func (s *MySQLSchema) GetTableName() string {
    return "whatever"
}

// AutoIncrementField 自增字段名称，如果没有则返回空
func (s *MySQLSchema) AutoIncrementField() string {
    return "whatever"
}

// GetDBFieldTag 获取数据库字段映射tag
func (s *MySQLSchema) GetDBFieldTag() string {
    return "whatever"
}

// MySQLSchemaModel model for MySQLSchema
type MySQLSchemaModel struct {
    *gomodel.ModelManager
}

// NewMySQLSchemaModel create a MySQLSchema Model
func NewMySQLSchemaModel(dsn string) *MySQLSchemaModel {
    m := &MySQLSchemaModel{
        msMgr,
    }
    m.SetDBInitFunc(func() (db *sql.DB, e error) {
        conn, err := sql.Open("mysql", dsn)
        if err != nil {
            return nil, err
        }
        return conn, nil
    })
    return m
}

// 查询数据表列表
func (m *MySQLSchemaModel) GetTables(database string) ([]map[string]string, error) {
    return m.NewQuerier().
        From("TABLES").
        Where(map[string]interface{}{"TABLE_SCHEMA": database}).
        QueryAll()
}

// 查询表信息
func (m *MySQLSchemaModel) GetTable(database, tableName string) (map[string]string, error) {
    return m.NewQuerier().
        From("TABLES").
        Where(map[string]interface{}{"TABLE_SCHEMA": database, "TABLE_NAME": tableName}).
        QueryRow()
}

// 查询数据表字段列表
func (m *MySQLSchemaModel) GetTableColumns(database, tableName string) ([]map[string]string, error) {
    return m.NewQuerier().
        From("COLUMNS").
        Where(map[string]interface{}{"TABLE_SCHEMA": database, "TABLE_NAME": tableName}).
        OrderBy("ORDINAL_POSITION ASC").
        QueryAll()
}

//-------------------- MODEL DEFINITION ---------------------//

// mysql gomodel生成器
type MySQLGenerator struct {
    conf *config.ItemConf
}

func NewMySQLGenerator(c *config.ItemConf) *MySQLGenerator {
    return &MySQLGenerator{
        conf: c,
    }
}

func (g *MySQLGenerator) checkConfig() error {
    if g.conf == nil {
        return fmt.Errorf("empty config")
    }
    if g.conf.OutputDir == "" || !fileutil.Exists(g.conf.OutputDir) {
        return fmt.Errorf("output dir not set or empty")
    }
    if g.conf.Database == "" {
        return fmt.Errorf("database not set or empty")
    }
    return nil
}

// 生成model文件 - 生成全部数据库的model文件
func (g *MySQLGenerator) Make() (bool, error) {
    // 配置检查
    err := g.checkConfig()
    if err != nil {
        return false, err
    }
    // 创建model
    m := NewMySQLSchemaModel(g.conf.DSN)
    // 查询全部数据表列表
    tables, err := m.GetTables(g.conf.Database)
    if err != nil {
        return false, fmt.Errorf("query tables of database %s failed: %s", g.conf.Database, err)
    }
    errCnt := 0
    var lastErr error = nil
    for _, tableInfo := range tables {
        _, err := g.makeByTableInfo(m, tableInfo)
        if err != nil {
            errCnt++
            lastErr = err
            fmt.Printf("genetate %s model file failed: %s", tableInfo["TABLE_NAME"], err)
        }
    }
    return errCnt == 0, lastErr
}

// 生成model文件 - 生成指定数据表的model文件
func (g *MySQLGenerator) MakeTable(table string) (bool, error) {
    // 配置检查
    err := g.checkConfig()
    if err != nil {
        return false, err
    }
    // 创建model
    m := NewMySQLSchemaModel(g.conf.DSN)
    tableInfo, err := m.GetTable(g.conf.Database, table)
    if err != nil {
        return false, err
    }
    return g.makeByTableInfo(m, tableInfo)
}

// 根据数据表信息生成model文件
func (g *MySQLGenerator) makeByTableInfo(m *MySQLSchemaModel, tableInfo map[string]string) (bool, error) {
    if tableInfo == nil || len(tableInfo) == 0 {
        return false, fmt.Errorf("empty table info")
    }
    tableName := tableInfo["TABLE_NAME"]
    tableComment := tableInfo["TABLE_COMMENT"]
    fmt.Printf("create model for table %s\r\n", tableName)
    // 文件基本信息转换
    fm := &filemaker.FileMaker{
        OutputDir:       g.conf.OutputDir,
        FileName:        fmt.Sprintf("%s.go", tableName),
        PackageName:     g.conf.Package,
        TableName:       tableName,
        TableComment:    tableComment,
        ConnectionName:  g.conf.ConnectionName,
        PrimaryKeyField: "",
        Fields:          make([]filemaker.Field, 0),
    }
    if fm.ConnectionName == "" {
        fm.ConnectionName = g.conf.Database
    }
    fmt.Printf("-- output dir: %s\r\n", fm.OutputDir)
    fmt.Printf("-- output file: %s\r\n", fm.FileName)
    // 字段处理
    // 查询字段列表
    columns, err := m.GetTableColumns(g.conf.Database, tableName)
    if err != nil {
        return false, err
    }
    if columns == nil || len(columns) == 0 {
        return false, fmt.Errorf("query table [%s] columns failed", tableName)
    }
    for _, column := range columns {
        field := filemaker.Field{}
        field.FieldName = column["COLUMN_NAME"]
        field.PropName = filemaker.GetPropName(field.FieldName)
        field.DataType = g.convertDataType(column["DATA_TYPE"], column["COLUMN_TYPE"])
        field.Comment = column["COLUMN_COMMENT"]
        fm.Fields = append(fm.Fields, field)
        // 检查是否是主键
        if column["COLUMN_KEY"] == "PRI" {
            fm.PrimaryKeyField = column["COLUMN_NAME"]
        }
    }
    // 生成文件
    return fm.MakeFile()
}

// 数据类型转换
func (g *MySQLGenerator) convertDataType(dataType, columnType string) string {
    t := strings.ToLower(strings.TrimSpace(dataType))
    isUnsigned := false
    if strings.Contains(columnType, "unsigned") {
        isUnsigned = true
    }
    if t == "bigint" || t == "int" {
        if isUnsigned {
            return "uint64"
        } else {
            return "int64"
        }
    }
    if t == "tinyint" || t == "smallint" || t == "bit" {
        if isUnsigned {
            return "uint"
        } else {
            return "int"
        }
    }
    if t == "decimal" || t == "float" {
        return "float64"
    }
    // 默认都返回字符串
    return "string"
}
