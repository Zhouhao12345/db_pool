package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm_demo/library"
	"log"
	"reflect"
	"strconv"
	"sync"
)

type Product struct {
	gorm.Model
	Code string
	Price uint
}


func testDBPool(index int, wg sync.WaitGroup) {
	defer wg.Done()
	log.Print(index)
	for i:=0; i < 15; i++ {
		cursor, err := library.Pool.Borrow()
		if err != nil {
			panic(err)
		}
		log.Printf("logger %d, pool length %d, flag %d, %s",
			cursor.Logger, library.Pool.Pool.Len(), index * i, strconv.FormatBool(cursor.Used))
		cursor.DB.Create(&Product{Code: "L1212", Price: uint(i)})
		library.Pool.Back(cursor)
	}
	return
}

func SingleUpdate(c *library.Cursor)  {
	c.DB.Create(&Product{Code: "L1212", Price: 100})
}

func main() {
	//err := library.Pool.DBContext(SingleUpdate)
	//if err != nil {
	//	panic(err)
	//}
	model := reflect.ValueOf(&Product{Code: "L1212", Price: 100})
	fmt.Println("Value", model)
	//var wg sync.WaitGroup
	//wg.Add(20)
	//for i:=0; i<100; i++ {
	//	go testDBPool(i, wg)
	//}
	//wg.Wait()
	//q := &util.Queue{}
	//for i:=0; i<10; i++ {
	//	q.Append(i)
	//	fmt.Print(q.Length)
	//}
	//for {
	//	n, _ := q.Pop()
	//	fmt.Print(n.Value)
	//}
}