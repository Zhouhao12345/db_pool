package configs

import (
	"gorm_demo/library/types"
)

var DbConfig = new(Config)

func init() {
	var dbValues types.ConfigMeta = map[string]interface{}{
		"username": "root",
		"password": "gllue123",
		"dbName":   "gllueweb_demo",
		"dbHost":   "0.0.0.0",
		"dbPort":   "3306",
	}
	DbConfig = &Config{
		Path:   "local",
		Alias:  "database",
		Values: dbValues,
	}
}
