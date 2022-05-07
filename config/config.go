package config

import (
    "github.com/BurntSushi/toml"
)

// 定义配置
type ItemConf struct {
    Driver         string `toml:"driver"`     // 数据库类型
    DSN            string `toml:"dsn"`        // 数据库连接
    Database       string `toml:"database"`   // 数据库名称
    ConnectionName string `toml:"conn_name"`  // 数据库连接名称
    OutputDir      string `toml:"output_dir"` // 输出目录
    Package        string `toml:"package"`    // 包名
}

type Config struct {
    DbItems []*ItemConf            `toml:"item"`
}

// 加载配置
// f - 配置文件
func LoadItems(f string) ([]*ItemConf, error) {
    cfg := &Config{}
    _, err := toml.DecodeFile(f, cfg)
    if err != nil {
        return nil, err
    }
    return cfg.DbItems, nil
}
