package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DataFile describes the simple structure of found .kt-file,
// that consists the filename and raw code as a content
type DataFile struct {
	Name    string
	Content []byte
}

type DataFiles struct {
	Files []DataFile
}

// Scan function takes the folder name, and recursively goes through the directory, and also its subdirectories
// return a list of all found kotlin files
func Scan(folder string) (DataFiles, error) {
	var files, err = RecursiveScanFolder(folder)
	return files, err
}

func RecursiveScanFolder(folder string) (DataFiles, error) {
	result := DataFiles{}
	return ScanFolders(&result, folder)
}

func ScanFolders(result *DataFiles, folder string) (DataFiles, error) {
	f, err := os.Open(folder)
	if err != nil {
		return *result, fmt.Errorf("the folder %s cannot be opened", folder)
	}

	files, err := f.Readdir(-1)
	_ = f.Close()
	if err != nil {
		return *result, fmt.Errorf("the folder %s cannot be read", folder)
	}

	for _, file := range files {
		if file.IsDir() {
			*result, err = ScanFolders(result, filepath.Join(folder, file.Name()))
		}
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".kt") {
			contentBytes, readErr := os.ReadFile(folder + "/" + file.Name())
			if readErr != nil {
				return *result, fmt.Errorf("the file %s cannot be read", file.Name())
			}
			result.Files = append(result.Files, DataFile{
				Name:    file.Name(),
				Content: contentBytes,
			})
		}
	}
	return *result, err
}
