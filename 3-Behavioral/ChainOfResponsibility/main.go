package main

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(*Context)

type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	index    int
	handlers []HandlerFunc
}

func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}

type Engine struct {
	handlers []HandlerFunc
}

func (e *Engine) Use(middleware HandlerFunc) {
	e.handlers = append(e.handlers, middleware)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		Request:  r,
		Writer:   w,
		index:    -1,
		handlers: e.handlers,
	}
	c.Next()
}

func LoggingMiddleware(c *Context) {
	fmt.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
	c.Next()
}

func AuthMiddleware(c *Context) {
	token := c.Request.Header.Get("Authorization")
	if token != "valid-token" {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("401 Unauthorized\n"))
		c.Abort()
		return
	}
	c.Next()
}

func HelloHandler(c *Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("Hello, World!\n"))
}

func main() {
	e := &Engine{}

	e.Use(LoggingMiddleware)
	e.Use(AuthMiddleware)
	e.Use(HelloHandler)

	fmt.Println("Server is starting...")
	http.ListenAndServe(":8080", e)
}
