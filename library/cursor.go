package library

import (
	"fmt"
	"gorm_demo/library/configs"
	"gorm_demo/library/types"
	"gorm_demo/util"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Cursor struct {
	DB   *gorm.DB
	Used bool
	Con  *configs.Config
}

func DBConnect(c *configs.Config) (cur *Cursor, err error) {
	cur = &Cursor{
		Used: false,
		Con:  c,
	}
	err = cur.New()
	return
}

func (c *Cursor) New() (err error) {
	username := c.Con.Get("username").(string)
	password := c.Con.Get("password").(string)
	dbName := c.Con.Get("dbName").(string)
	dbHost := c.Con.Get("dbHost").(string)
	dbPort := c.Con.Get("dbPort").(string)
	dbConnectURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, dbHost, dbPort, dbName,
	)
	db, err := gorm.Open(
		"mysql",
		dbConnectURL,
	)
	if err != nil {
		return
	}
	c.DB = db
	return
}

func (c *Cursor) SetMaxConnect(number int) (err error) {
	db := c.DB.DB()
	db.SetMaxOpenConns(number)
	return
}

func (c *Cursor) Close() (err error) {
	err = c.DB.Close()
	return
}

func (c *Cursor) GetUsed() (used bool) {
	return c.Used
}

func (c *Cursor) SetUsed(used bool) (err error) {
	c.Used = used
	return
}

func (c *Cursor) HandlerRequest(f types.HandlerFunc) types.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("db", c)
		f(context)
		return
	}
}

func (c *Cursor) Migrate(value ...interface{}) {
	c.DB.AutoMigrate(value...)
}

type DatabaseManager struct {
	ConnectManager
}

var DBPool *DatabaseManager

func init() {
	NewConnectHandler := func() (types.Connect, error) {
		var (
			c   types.Connect
			err error
		)
		c, err = DBConnect(configs.DbConfig)
		return c, err
	}
	locker := sync.Mutex{}
	cond := sync.NewCond(&locker)
	DBPool = &DatabaseManager{
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
