package library

import (
	"gorm_demo/library/configs"
	"gorm_demo/library/types"
	"gorm_demo/util"
	"sync"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
)

type Cache struct {
	Client *redis.Client
	Used   bool
	Con    *configs.Config
}

func CacheConnect(c *configs.Config) (cli *Cache, err error) {
	cli = &Cache{
		Used: false,
		Con:  c,
	}
	err = cli.New()
	return
}

func (c *Cache) New() (err error) {
	host := c.Con.Get("host").(string)
	password := c.Con.Get("password").(string)
	db := c.Con.Get("db").(int)
	client := redis.NewClient(
		&redis.Options{
			Addr:     host,
			Password: password,
			DB:       db,
		})
	c.Client = client
	return
}

func (c *Cache) Close() (err error) {
	err = c.Client.Close()
	return
}

func (c *Cache) GetUsed() (used bool) {
	return c.Used
}

func (c *Cache) SetUsed(used bool) (err error) {
	c.Used = used
	return
}

func (c *Cache) HandlerRequest(f types.HandlerFunc) types.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("cache", c)
		f(context)
		return
	}
}

type CacheManager struct {
	ConnectManager
}

var CachePool *CacheManager

func init() {
	NewConnectHandler := func() (types.Connect, error) {
		var (
			c   types.Connect
			err error
		)
		c, err = CacheConnect(configs.RedisConfig)
		return c, err
	}
	locker := sync.Mutex{}
	cond := sync.NewCond(&locker)
	CachePool = &CacheManager{
		ConnectManager{
			Pool: &PoolManager{
				MinConnect: 10,
				MaxConnect: 5000,
				Pool:       new(util.Queue),
				lock:       cond,
				New:        NewConnectHandler,
			},
		},
	}
	DBPool.New()
}
