package main

import "fmt"

type Req interface {
	Handler(url string) string
}

type Request struct{}

func (r *Request) Handler(url string) string {
	return "Request to " + url
}

func main() {
	req := &LogDecorator{req: &Request{}}
	println(req.Handler("www.example.com"))
}

type LogDecorator struct {
	req Req
}

func (l *LogDecorator) Handler(url string) string {
	fmt.Println("Logging ...")
	return l.req.Handler(url)
}
