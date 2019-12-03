package main

import (
	"fmt"
	_ "fmt"
	"io"
	_ "io"
	"io/ioutil"
	"os"
	_ "path/filepath"
	"strconv"
	_ "strings"
)

const testDirResult = `├───project
├───static
│	├───a_lorem
│	│	└───ipsum
│	├───css
│	├───html
│	├───js
│	└───z_lorem
│		└───ipsum
└───zline
	└───lorem
		└───ipsum
`

const testFullResult = `├───project
│	├───file.txt (19b)
│	└───gopher.png (70372b)
├───static
│	├───a_lorem
│	│	├───dolor.txt (empty)
│	│	├───gopher.png (70372b)
│	│	└───ipsum
│	│		└───gopher.png (70372b)
│	├───css
│	│	└───body.css (28b)
│	├───empty.txt (empty)
│	├───html
│	│	└───index.html (57b)
│	├───js
│	│	└───site.js (10b)
│	└───z_lorem
│		├───dolor.txt (empty)
│		├───gopher.png (70372b)
│		└───ipsum
│			└───gopher.png (70372b)
├───zline
│	├───empty.txt (empty)
│	└───lorem
│		├───dolor.txt (empty)
│		├───gopher.png (70372b)
│		└───ipsum
│			└───gopher.png (70372b)
└───zzfile.txt (empty)
`

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(buffer io.Writer, s string, b bool) error {

	strings := make([]string, 0)
	prefix := "├───"
	splitter := ""
	err := WalkString(s+string(os.PathSeparator)+"testdata", &strings, prefix, splitter, b, false)

	var mainString string
	for _, s := range strings  {
		mainString += s + "\n"
	}

	if b {
		fmt.Println(mainString)
	}

	fmt.Println(mainString == testDirResult)
	fmt.Println(mainString == testFullResult)

	if err != nil {
		return err
	}
	return err
}

func WalkString(s string, strings *[]string, prefix string, splitter string, b bool, last bool) error {

	dir, err := ioutil.ReadDir(s)
	if !b {
		dir = justDir(dir)
	}

	if err != nil {
		return err
	}

	for index, f := range dir {
		if index == len(dir)-1 {
			prefix = "└───"
		} else {
			prefix = "├───"
		}
		var name string

		name = splitter + prefix + f.Name()
		if b && !f.IsDir() {
			name += getFileSize(f)
		}
		if len(splitter) != 0 {
			name = "│" + name
		}

		if len(splitter) == 0 && index == len(dir)-1 {
			last = true
		}

		if last {
			name = splitter + prefix + f.Name()
			if b && !f.IsDir() {
				name += getFileSize(f)
			}
		}
		//if b {
		//	name += getFileSize(f)
		//}

		*strings = append(*strings, name)
		//fmt.Println(name)
		if f.IsDir() {
			if len(splitter) == 0 {
				err = WalkString(s+string(os.PathSeparator)+f.Name(), strings, prefix, splitter+"\t", b, last)
			} else if index == len(dir)-1 {
				err = WalkString(s+string(os.PathSeparator)+f.Name(), strings, prefix, splitter+"\t", b, last)
			} else {
				err = WalkString(s+string(os.PathSeparator)+f.Name(), strings, prefix, splitter+"│"+"\t", b, last)
			}
			if err != nil {
				return err
			}
		}
	}

	return err
}

func justDir(files []os.FileInfo) []os.FileInfo {
	var list []os.FileInfo
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		list = append(list, f)
	}
	return list
}

func getFileSize(f os.FileInfo) string  {
	size := f.Size()
	var name string
	if size != 0 {
		name += " (" + strconv.Itoa(int(size)) + "b)"
	} else {
		name += " " + "(empty)"
	}

	return name
}