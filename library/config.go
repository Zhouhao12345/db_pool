package library

type Config struct {
	Path string
	Alias string
	Values ConfigMeta
}

func (c *Config) Get(key string) interface{} {
	return c.Values[key]
}
