package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:1900271083@tcp(43.138.107.154:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	if err := InitDB("mysql", dsn, 10, 5, time.Minute); err != nil {
		panic(err)
	}
	defer CloseDB()

	db := GetDB()
	StmtCache := NewStmtCache(db)
	defer StmtCache.Close()

	repo, err := NewUserRepo(db, StmtCache)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	RegisterHandlers(mux, repo)
	go func() {
		fmt.Println("HTTP listening on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			panic(err)
		}
	}()

	// 这里可以执行插入测试用户等启动后任务
	// 插入测试用户
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := &User{
		Name:      "Test",
		Email:     "1900271083@qq.com",
		CreatedAt: time.Now(),
	}

	id, err := repo.Create(ctx, u)
	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			fmt.Println("测试用户已存在")
		} else {
			panic(err)
		}
	}
	fmt.Println("Inserted user ID:", id)

	select {} // 阻塞主协程，或用信号优雅退出
}
