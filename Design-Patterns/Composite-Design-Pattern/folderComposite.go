package main

import (
	"fmt"
	"strings"
)

type FolderComposite struct{
	name string
	folder []FileComponent
}

func NewFolderComposite(name string)*FolderComposite{
	return &FolderComposite{name,make([]FileComponent, 0)}
}

func(f *FolderComposite)ShowDetails(){
	fmt.Println("+ ",f.name)
	for _,fol:=range f.folder{
		fol.ShowDetails()
	}

}

func(f *FolderComposite)GetSize()int{
	val:=0
	for _,v:=range f.folder{
		val+=v.GetSize()
	}
	return val

}

func(f *FolderComposite)Search(name string)bool{
	name=strings.ToLower(name)
	if name==f.name{
		return true
	}
	for _,val:=range f.folder{
		if val.Search(name){
			return true
		}
	}
	return false
}

func(f *FolderComposite)AddFileSystem(file FileComponent){
	f.folder = append(f.folder, file)
}