package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type KotlinFile struct {
	Name    string
	Content []byte
}

type KotlinFiles struct {
	Files []KotlinFile
}

func Scan(folder string) (KotlinFiles, error) {
	var files, err = RecursiveScanFolder(folder)
	return files, err
}

func RecursiveScanFolder(folder string) (KotlinFiles, error) {
	result := KotlinFiles{}
	return ScanFolders(&result, folder)
}

func ScanFolders(result *KotlinFiles, folder string) (KotlinFiles, error) {
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
			result.Files = append(result.Files, KotlinFile{
				Name:    file.Name(),
				Content: contentBytes,
			})
		}
	}
	return *result, err
}
