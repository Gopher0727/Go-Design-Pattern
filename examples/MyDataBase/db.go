package main

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"
)

var (
	db    *sql.DB
	once  sync.Once
	dbErr error
)

func InitDB(driver, dsn string, maxOpen, maxIdle int, connMaxLifetime time.Duration) error {
	once.Do(func() {
		db, dbErr = newDB(driver, dsn)
		if dbErr != nil {
			return
		}
		db.SetMaxOpenConns(maxOpen)
		db.SetMaxIdleConns(maxIdle)
		db.SetConnMaxLifetime(connMaxLifetime)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		dbErr = db.PingContext(ctx)
		if dbErr != nil {
			_ = db.Close()
			return
		}
	})
	return dbErr
}

func newDB(driver, dsn string) (*sql.DB, error) {
	switch driver {
	case "mysql":
		return sql.Open("mysql", dsn)
	case "postgres":
		return sql.Open("postgres", dsn)
	default:
		return nil, errors.New("unsupported driver: " + driver)
	}
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func HealthCheck(ctx context.Context) error {
	if db == nil {
		return sql.ErrConnDone
	}
	return db.PingContext(ctx)
}
