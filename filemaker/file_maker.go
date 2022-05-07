package filemaker

import (
    "fmt"
    "os/exec"
    "strings"
    "unicode"
)

// 文件生成器
type FileMaker struct {
    OutputDir       string  // 输出目录
    FileName        string  // 文件名称
    PackageName     string  // 包名
    TableName       string  // 数据表名
    TableComment    string  // 注释
    ConnectionName  string  // 配置的数据库连接名
    PrimaryKeyField string  // 主键字段，只支持一个字段，不支持联合主键
    Fields          []Field // 数据字段列表
}

// 定义表格字段
type Field struct {
    FieldName string // 数据库原生字段名称
    PropName  string // 生成目标文件属性名称
    DataType  string // 数据类型名称
    Comment   string // 注释信息
}

// Ucfirst ucfirst()
func Ucfirst(str string) string {
    for _, v := range str {
        u := string(unicode.ToUpper(v))
        return u + str[len(u):]
    }
    return ""
}

// 将字段转换为go结构字段
func GetPropName(s string) string {
    s = strings.ToLower(s)
    if s == "id" {
        return "ID"
    }
    parts := strings.Split(s, "_")
    retS := ""
    for _, p := range parts {
        if p == "id" {
            retS += "ID"
            continue
        }
        retS += Ucfirst(p)
    }
    return retS
}

// 将ddl关键字段转换为go结构字段
func GetVarName(s string) string {
    s = strings.ToLower(s)
    if s == "id" {
        return "ID"
    }
    parts := strings.Split(s, "_")
    retS := ""
    for i, p := range parts {
        if i == 0 {
            retS += p
            continue
        }
        if p == "id" {
            retS += "ID"
            continue
        }
        retS += Ucfirst(p)
    }
    return retS
}

// 格式化go文件
func FmtGoFile(f string) {
    fmt.Printf("-- format file: %s \r\n", f)
    cmd := exec.Command("gofmt", "-l", "-w", "-s", f)
    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("-- failed: %s\n", string(out))
        fmt.Println("----------------")
        fmt.Printf("-- gofmt failed: %s\n", err)
        return
    }
    fmt.Printf("-- success: %s\n", string(out))
}
