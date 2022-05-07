package generator

import (
    "fmt"
    "gomodelcreator/config"
)

// 默认的生成器
type DefaultGenerator struct {
    conf *config.ItemConf
}

func NewDefaultGenerator(c *config.ItemConf) *DefaultGenerator {
    return &DefaultGenerator{
        conf: c,
    }
}

func (g *DefaultGenerator) Make() (bool, error) {
    return false, fmt.Errorf("[%s] generator not supported", g.conf.Driver)
}

func (g *DefaultGenerator) MakeTable(table string) (bool, error) {
    return false, fmt.Errorf("[%s] generator not supported", g.conf.Driver)
}
