package configs

import (
	"gorm_demo/library/types"
)

var RedisConfig = new(Config)

func init() {
	var cacheValues types.ConfigMeta = map[string]interface{}{
		"host":     "0.0.0.0:6379",
		"password": "",
		"db":       0,
	}
	RedisConfig = &Config{
		Path:   "local",
		Alias:  "cache",
		Values: cacheValues,
	}
}
