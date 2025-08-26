package main

import (
	"fmt"
	"time"
)

type Image interface {
	Display()
}

type RealIamge struct {
	filename string
}

func NewRealImage(filename string) *RealIamge {
	r := &RealIamge{filename: filename}
	r.loadFromDisk()
	return r
}

func (r *RealIamge) loadFromDisk() {
	fmt.Println("Loading image from disk...")
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Image loaded. %s\n", r.filename)
}

func (r *RealIamge) Display() {
	fmt.Printf("Displaying %s\n", r.filename)
}

type ProxyImage struct {
	filename  string
	realImage *RealIamge
}

func NewProxyImage(filename string) *ProxyImage {
	return &ProxyImage{filename: filename}
}

func (p *ProxyImage) Display() {
	if p.realImage == nil {
		fmt.Println("Proxy: Creating real image...")
		p.realImage = NewRealImage(p.filename)
	} else {
		fmt.Println("Proxy: Real image already created.")
	}
	// do something
	p.realImage.Display()
}

func main() {
	img := NewProxyImage("test_image.jpg")

	img.Display()
	img.Display()
}
