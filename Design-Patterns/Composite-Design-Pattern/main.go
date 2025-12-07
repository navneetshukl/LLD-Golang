package main

import "fmt"

func main() {
	file1 := NewFileLeaf("file1.txt", 10)
	file2 := NewFileLeaf("file2.txt", 15)
	file3 := NewFileLeaf("file.json", 8)

	folder1 := NewFolderComposite("folder1")
	folder2 := NewFolderComposite("folder2")
	folder3 := NewFolderComposite("folder3")

	file4 := NewFileLeaf("file4.html", 10)
	file5 := NewFileLeaf("file5.java", 15)
	file6 := NewFileLeaf("file6.go", 8)

	file7 := NewFileLeaf("file7.json", 10)
	file8 := NewFileLeaf("file8.png", 15)
	file9 := NewFileLeaf("file9.jpeg", 8)

	folder1.AddFileSystem(file1)
	folder1.AddFileSystem(file2)
	folder1.AddFileSystem(file5)
	folder2.AddFileSystem(file6)
	folder2.AddFileSystem(folder1)

	folder2.AddFileSystem(file3)
	folder2.AddFileSystem(file4)

	folder3.AddFileSystem(file7)
	folder3.AddFileSystem(file8)
	folder3.AddFileSystem(file9)

	folder3.AddFileSystem(folder2)

	folder3.ShowDetails()
	fmt.Println(folder2.Search(folder1.name))

}
