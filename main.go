package main

import (
    "flag"
    "fmt"
    "gomodelcreator/config"
    "gomodelcreator/generator"
)

func main() {
    // 指定配置文件
    configFile := flag.String("f", "config.toml", "config file")
    // 只生成指定表的model文件
    db := flag.String("d", "", "database name")
    table := flag.String("t", "", "table name")
    flag.Parse()
    // 初始化配置信息
    cfgs, err := config.LoadItems(*configFile)
    if err != nil {
        fmt.Printf("load config failed： %s\r\n", err)
        return
    }
    if len(cfgs) == 0 {
        fmt.Printf("config empty or format error\r\n")
        return
    }
    // 如果指定了表，则必须指定数据库
    if *table != "" {
        if *db == "" {
            fmt.Printf("database required when create model for a specified table\r\n")
            return
        }
    }
    // 根据配置生成model文件
    for _, conf := range cfgs {
        // 检查是全部刷新还是只刷新单个model
        if *db != "" {
            if conf.Database != *db {
                continue
            }
            // 检查是否只是刷新某个表的model
            if *table != "" {
                fmt.Printf("create model for table %s.%s...\r\n", *db, *table)
                // 刷新指定的model
                gen := generator.NewGenerator(conf)
                _, err = gen.MakeTable(*table)
                if err != nil {
                    fmt.Printf("create mode for table %s failed：%s\r\n", *table, err)
                }
            } else {
                fmt.Printf("create model for db %s...\r\n", *db)
                // 刷新指定的model
                gen := generator.NewGenerator(conf)
                _, err = gen.Make()
                if err != nil {
                    fmt.Printf("create mode for db %s failed (maybe partially)：%s\r\n", *db, err)
                }
            }
            continue
        }
        // 数据库中的表全量生成
        fmt.Printf("create model for db %s...\r\n", conf.Database)
        gen := generator.NewGenerator(conf)
        _, err = gen.Make()
        if err != nil {
            fmt.Printf("create mode for db %s failed (maybe partially)：%s\r\n", conf.Database, err)
        }
    }
}


