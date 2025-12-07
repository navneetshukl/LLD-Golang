package main

type FileComponent interface {
	ShowDetails()
	GetSize() int
	Search(name string) bool
}
