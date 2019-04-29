package configs

import (
	"gorm_demo/library/types"
)

type Config struct {
	Path   string
	Alias  string
	Values types.ConfigMeta
}

func (c *Config) Get(key string) interface{} {
	return c.Values[key]
}
