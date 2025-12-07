package main

import (
	"fmt"
	"strings"
)

type FileLeaf struct {
	name string
	size int
}

func NewFileLeaf(name string, size int) *FileLeaf {
	return &FileLeaf{name, size}
}

func (f *FileLeaf) ShowDetails() {
	fmt.Printf("Name : %s || Size : %d \n", f.name, f.size)

}

func (f *FileLeaf) GetSize() int {
	return f.size
}

func (f *FileLeaf) Search(name string) bool {
	name = strings.ToLower(name)
	return name == f.name
}
