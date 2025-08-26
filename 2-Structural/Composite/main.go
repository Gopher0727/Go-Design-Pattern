package main

import "fmt"

type Node interface {
	Display(indentation string)
}

type File struct {
	Name string
}

func (f *File) Display(indentation string) {
	fmt.Println(indentation + f.Name)
}

type Dir struct {
	Name  string
	Nodes []Node
}

func (d *Dir) Display(indentation string) {
	fmt.Println(indentation + d.Name)
	for _, node := range d.Nodes {
		node.Display(indentation + "  ")
	}
}

func main() {
	root := &Dir{
		Name: "root",
		Nodes: []Node{
			&Dir{
				Name:  "bin",
				Nodes: []Node{},
			},
			&Dir{
				Name: "usr",
				Nodes: []Node{
					&File{
						Name: "config.toml",
					},
				},
			},
			&File{
				Name: "files.txt",
			},
		},
	}
	root.Display("")
}
