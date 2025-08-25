package main

// 饿汉式
type DataBase struct {
	dsn string
}

var db *DataBase

func init() {
	db = &DataBase{
		dsn: "your_dsn_here",
	}
}

func GetDB() *DataBase {
	return db
}
