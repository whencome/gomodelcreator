package generator

import (
    "gomodelcreator/config"
)

// 定义一个生成器
type Generator interface {
    Make() (bool, error)
    MakeTable(table string) (bool, error)
}

// 工厂方法类，用于获取一个生成器
func NewGenerator(conf *config.ItemConf) Generator {
    switch conf.Driver {
    case "mysql":
        return NewMySQLGenerator(conf)
    default:
        return NewDefaultGenerator(conf)
    }
}