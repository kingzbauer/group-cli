package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var dirPath = flag.String("dir", "", "The base directory to process")

func ExtractOrdinaryFiles(basePath string, dirNames []string) []string {
	fileNames := make([]string, len(dirNames))
	numOfFiles := 0 // count of ordinary files

	// Construct path for individual files
	for _, name := range dirNames {
		filePath := filepath.Join(basePath, name)
		// Try to open file
		file, err := os.Open(filePath)
		if err != nil {
			continue
		}
		fileInfo, err := file.Stat()
		// if not ordinary file, skip
		if err != nil {
			continue
		}
		if fileInfo.IsDir() {
			continue
		}

		fileNames[numOfFiles] = filePath
		numOfFiles += 1
	}

	finalFileNames := make([]string, numOfFiles)
	copy(finalFileNames, fileNames)

	return finalFileNames
}

func GroupByExt(fileNames []string) (group map[string][]string) {
	group = make(map[string][]string)

	for _, path := range fileNames {
		ext := filepath.Ext(path)
		if len(ext) == 0 {
			ext = ".unknown"
		}
		// remove the leading dot
		ext = strings.ToLower(ext[1:])

		// add to map
		group[ext] = append(group[ext], path)
	}
	return
}

func MoveToDir(dirPath string, fileNames []string) {
	dir, err := os.Open(dirPath)
	// If err, no dir, create
	if err != nil {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error while creating dir:", err)
			return
		}
	}

	dir, err = os.Open(dirPath)
	// Directory info
	dirInfo, err := dir.Stat()
	if err != nil {
		// Unexpected error
		fmt.Println("Error while retrieving directory stat:", err)
		return
	}

	// If not a directory, create one
	if !dirInfo.IsDir() {
		err = os.Mkdir(dirPath, os.ModeDir)
	}

	// rename each file
	for _, oldPath := range fileNames {
		// get file name
		_, fileName := filepath.Split(oldPath)
		// new path
		newPath := filepath.Join(dirPath, fileName)
		os.Rename(oldPath, newPath)
	}
}

func main() {
	flag.Parse()
	basePath := *dirPath

	// If not directory name specified, return
	if len(basePath) == 0 {
		fmt.Println("Please specify a directory")
		return
	}

	// Open the directory
	dir, err := os.Open(basePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	dirInfo, err := dir.Stat()

	// If not directory, return
	if !dirInfo.IsDir() {
		fmt.Printf("%s is not a directory\n", dirInfo.Name())
		return
	}

	// Retrieve the directory names
	dirNames := []string{}
	dirNames, err = dir.Readdirnames(0)
	if err != nil {
		// Show the error info
		fmt.Println("Error:", err.Error())
	}
	if !(len(dirNames) > 0) {
		fmt.Println("This directory does not contain any files")
		return
	}

	// Extract ordinary files
	ordinaryFiles := ExtractOrdinaryFiles(basePath, dirNames)
	grouped := GroupByExt(ordinaryFiles)
	for dir, files := range grouped {
		dirPath := filepath.Join(basePath, dir)
		MoveToDir(dirPath, files)
	}
}
