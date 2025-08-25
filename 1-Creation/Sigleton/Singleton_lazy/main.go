package main

import (
	"sync"
)

// 懒汉式
type DataBase struct {
	dsn string
}

var (
	db   *DataBase
	once sync.Once
)

func GetDB() *DataBase {
	once.Do(func() {
		db = &DataBase{
			dsn: "your_dsn_here",
		}
	})
	return db
}
