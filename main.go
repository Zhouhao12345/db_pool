package main

import (
	"gorm_demo/library"
	"gorm_demo/library/configs"
	"gorm_demo/library/types"
	"gorm_demo/models"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

// func testDBPool(index int, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for i := 0; i < 15; i++ {
// 		library.DBPool.DBContext(Create)(&models.Product{Code: "L1212", Price: uint(i)})
// 	}
// }

func Create(context *types.Context) {
	cursor, existed := context.Get("db")
	if !existed {
		panic("No Existd DB Connect")
	}
	cursor.(*library.Cursor).DB.Create(
		&models.Product{Code: "L1212", Price: 100})
}

func Print(context *types.Context) {
	cache, existed := context.Get("cache")
	if !existed {
		panic("No Existd DB Connect")
	}
	_, err := cache.(*library.Cache).Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func BaseMigrate(db interface{}, value ...interface{}) (err error) {
	db.(*gorm.DB).AutoMigrate(value...)
	return
}

func test_db(wg *sync.WaitGroup) {
	defer wg.Done()
	context := new(types.Context)
	library.DBPool.Context(Create)(context)
}

func test_cache(wg *sync.WaitGroup) {
	defer wg.Done()
	context := new(types.Context)
	library.CachePool.Context(Print)(context)
}

func test_office(wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	db.Create(&models.Product{Code: "L1212", Price: 100})
}

func main() {
	start := time.Now()
	args := os.Args
	switch args[1] {
	case "cache":
		var wg sync.WaitGroup
		waitNum := 5000
		wg.Add(5000)
		for i := 0; i < waitNum; i++ {
			go test_cache(&wg)
		}
		wg.Wait()
	case "db":
		var wg sync.WaitGroup
		waitNum := 5000
		wg.Add(5000)
		for i := 0; i < waitNum; i++ {
			go test_db(&wg)
		}
		wg.Wait()
	case "run_office":
		var wg sync.WaitGroup
		waitNum := 5000
		wg.Add(5000)
		cur, _ := library.DBConnect(configs.DbConfig)
		cur.DB.DB().SetMaxOpenConns(5000)
		cur.DB.DB().SetMaxIdleConns(2)
		cur.DB.DB().SetConnMaxLifetime(time.Duration(14400) * time.Hour)
		for i := 0; i < waitNum; i++ {
			go test_office(&wg, cur.DB)
		}
		wg.Wait()
	}
	time_duration := time.Since(start)
	log.Println(time_duration)
}
