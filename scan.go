package main

import (
	"KotlinToJavaConverter/structures"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Scan function takes the folder name, and recursively goes through the directory, and also its subdirectories
// return a list of all found kotlin files
func Scan(folder string) (structures.DataFiles, error) {
	var files, err = RecursiveScanFolder(folder)
	return files, err
}

func RecursiveScanFolder(folder string) (structures.DataFiles, error) {
	result := structures.DataFiles{}
	return ScanFolders(&result, folder)
}

func ScanFolders(result *structures.DataFiles, folder string) (structures.DataFiles, error) {
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
			result.Files = append(result.Files, structures.DataFile{
				Name:    file.Name(),
				Content: contentBytes,
			})
		}
	}
	return *result, err
}
