package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

const (
	EmptyPostfix = " (empty)"
	LastElementSymbol = "└───"
	RegularElementSymbol = "├───"
	UniqueStick = "│"
)

type FileFormat struct {
	formatPrefix string
	fileInfo fs.FileInfo
	isFileLast bool
}

func (f *FileFormat) formatFileOutput(out io.Writer) {
	format := "%s"
	if !f.fileInfo.IsDir() {
		if f.fileInfo.Size() != 0 {
			format = format + " (" + strconv.Itoa(int(f.fileInfo.Size())) +"b)"
		} else {
			format = format + EmptyPostfix
		}
	}
	if f.isFileLast {
		format = f.formatPrefix + LastElementSymbol + format + "\n"
	} else {
		format = f.formatPrefix + RegularElementSymbol + format + "\n"
	}
	fmt.Fprintf(out, format, f.fileInfo.Name())
}

func inside(prefix string, out io.Writer, path string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name()})
	fileIterator := 0
	for _, val := range files {
		file := &FileFormat{
			formatPrefix: prefix,
			fileInfo:     val,
			isFileLast:   filesSum(files, printFiles) <= fileIterator + 1,
		}
		if printFiles || val.IsDir() {
			file.formatFileOutput(out)
			fileIterator++
		}
		if val.IsDir() {
			insidePrefix := prefix
			if !file.isFileLast {
				 insidePrefix += UniqueStick
			}
			err = inside(insidePrefix + "\t", out, path + "/" + file.fileInfo.Name(), printFiles)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func filesSum(files []fs.FileInfo, printFiles bool) int {
	lenDir := 0
	if printFiles {
		return len(files)
	}
	for _, val := range files {
		if val.IsDir() {
			lenDir++
		}
	}
	return lenDir
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := inside("", out, path, printFiles)
	if err != nil {
		return err
	}
	return nil
}

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
