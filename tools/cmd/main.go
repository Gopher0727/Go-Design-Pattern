package main

import (
	"flag"
	"fmt"
	"os"

	tree "small-tree/internal"
)

func main() {
	opts := tree.ParseFlags()
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	t, err := tree.Build(root, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := tree.Render(t, os.Stdout, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
