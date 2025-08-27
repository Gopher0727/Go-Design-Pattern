package main

import (
	"context"
	"database/sql"
	"sync"
)

type StmtCache struct {
	db    *sql.DB
	mu    sync.RWMutex
	cache map[string]*sql.Stmt
}

func NewStmtCache(db *sql.DB) *StmtCache {
	return &StmtCache{
		db:    db,
		cache: make(map[string]*sql.Stmt),
	}
}

func (c *StmtCache) Prepare(ctx context.Context, key, query string) (*sql.Stmt, error) {
	c.mu.RLock()
	if s, ok := c.cache[key]; ok {
		c.mu.RUnlock()
		return s, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if s, ok := c.cache[key]; ok {
		return s, nil
	}

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	c.cache[key] = stmt
	return stmt, nil
}

func (c *StmtCache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var firstErr error
	for key, s := range c.cache {
		if err := s.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(c.cache, key)
	}
	return firstErr
}
