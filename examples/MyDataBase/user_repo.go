package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, u *User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Close() error
}

type userRepo struct {
	db    *sql.DB
	cache *StmtCache
}

func NewUserRepo(db *sql.DB, cache *StmtCache) (UserRepository, error) {
	// 简单迁移
	schema := `CREATE TABLE IF NOT EXISTS users (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        created_at DATETIME NOT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &userRepo{db: db, cache: cache}, nil
}

func (r *userRepo) Close() error {
	return nil
}

func (r *userRepo) Create(ctx context.Context, u *User) (int64, error) {
	stmt, err := r.cache.Prepare(ctx, "insertUser", "INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.ExecContext(ctx, u.Name, u.Email, u.CreatedAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	stmt, err := r.cache.Prepare(ctx, "user.select_by_id", "SELECT id, name, email, created_at FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, id)
	u := &User{}
	var created time.Time
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	u.CreatedAt = created
	return u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	stmt, err := r.cache.Prepare(ctx, "user.select_by_email", "SELECT id, name, email, created_at FROM users WHERE email = ?")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, email)
	u := &User{}
	var created time.Time
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	u.CreatedAt = created
	return u, nil
}
