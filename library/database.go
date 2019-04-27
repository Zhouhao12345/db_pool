package library

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm_demo/util"
	"sync"
)

type Cursor struct {
	DB *gorm.DB
	Used bool
	Logger int
}

func DBConnect(c *Config) (cu *Cursor , err error) {
	username := c.Get("username").(string)
	password := c.Get("password").(string)
	dbName := c.Get("dbName").(string)
	dbHost := c.Get("dbHost").(string)
	dbPort := c.Get("dbPort").(string)
	dbConnectUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, dbHost, dbPort, dbName,
	)
	db, err := gorm.Open(
		"mysql",
		dbConnectUrl,
	)
	cu = &Cursor{
		DB: db,
		Used: false,
	}
	if err != nil {
		return
	}
	return
}

func (cur *Cursor) Migrate(value interface{}) {
	cur.DB.AutoMigrate(value)
}

type PoolManager struct {
	Pool *util.Queue
	All []*Cursor
	lock *sync.Cond
	BaseConfig *Config
	MinConnect int
	MaxConnect int
}

func (p *PoolManager) init()  {
	for i:=0; i<p.MinConnect; i++ {
		if cur, err := DBConnect(p.BaseConfig); err == nil {
			p.All = append(p.All, cur)
			p.Pool.Append(cur)
		} else {
			panic(err)
		}
	}
}

func (p *PoolManager) Borrow() (cur *Cursor, err error) {
	p.lock.L.Lock()
	for {
		if p.Pool.End == nil {
			if len(p.All) >= p.MaxConnect {
				p.lock.Wait()
			} else {
				cur, err = DBConnect(p.BaseConfig)
				if err != nil {
					panic(err)
				}
				cur.Used = true
				p.All = append(p.All, cur)
				p.lock.L.Unlock()
				return
			}
		} else {
			break
		}
	}

	var qn *util.QNode
	qn, err = p.Pool.Pop()
	if err != nil {
		return nil, err
	}
	cur = qn.Value.(*Cursor)

	cur.Used = true
	p.lock.L.Unlock()
	return
}

func (p *PoolManager) Back(cursor *Cursor) {
	p.lock.L.Lock()
	cursor.Used = false
	p.Pool.Append(cursor)
	p.lock.Broadcast()
	p.lock.L.Unlock()
	return
}

func (p *PoolManager) DBContext(f DBHandler) (err error) {
	cursor, err := p.Borrow()
	if err != nil {
		return
	}
	f(cursor)
	p.Back(cursor)
	return nil
}

var Pool = new(PoolManager)

func init()  {
	var dbValues ConfigMeta = map[string]interface{}{
		"username": "root",
		"password": "Hello.123",
		"dbName": "go_ws",
		"dbHost": "localhost",
		"dbPort": "3306",
	}
	dbConfig := &Config{
		Path:   "local",
		Alias:  "database",
		Values: dbValues,
	}
	locker := sync.Mutex{}
	cond := sync.NewCond(&locker)
	Pool = &PoolManager{
		MinConnect:5,
		MaxConnect:10,
		BaseConfig:dbConfig,
		Pool:new(util.Queue),
		lock:cond,
	}
	Pool.init()
}
