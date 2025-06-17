package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Helper Functions for Testing ---

// createTestDir creates a temporary directory with a specified structure of files and content.
// It returns the path to the created temporary directory.
// The `files` map uses relative paths within the temporary directory as keys
// and file contents as values.
func createTestDir(t *testing.T, files map[string]string) string {
	// Create a unique temporary directory for this test run
	tmpDir, err := os.MkdirTemp("", "kotlin_scanner_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Populate the temporary directory with specified files
	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)

		// Create parent directories if they don't exist
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tmpDir) // Clean up on error
			t.Fatalf("Failed to create parent dirs for %s: %v", fullPath, err)
		}

		// Write the file content
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			os.RemoveAll(tmpDir) // Clean up on error
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}
	return tmpDir
}

// compareKotlinFiles is a helper to compare two slices of DataFile structs.
// It normalizes paths to be relative to the test directory for consistent comparison
// and ignores the order of files.
func compareKotlinFiles(t *testing.T, got, want []DataFile, baseDir string) {
	// If lengths don't match, they are definitely different
	if len(got) != len(want) {
		t.Errorf("Expected %d files, got %d", len(want), len(got))
		return
	}

	// Create maps for easier comparison, using normalized file names/paths as keys
	gotMap := make(map[string]string)
	for _, f := range got {
		// Normalize path to be relative to baseDir for robust comparison
		relPath, err := filepath.Rel(baseDir, filepath.Join(baseDir, f.Name)) // Assuming f.Name is base name
		if err != nil {
			t.Fatalf("Error getting relative path for %s: %v", f.Name, err)
		}
		gotMap[relPath] = string(f.Content)
	}

	wantMap := make(map[string]string)
	for _, f := range want {
		wantMap[f.Name] = string(f.Content) // want.Name is already relative for test definition
	}

	// Check if all wanted files are present and have correct content in got
	for name, wantContent := range wantMap {
		gotContent, ok := gotMap[name]
		if !ok {
			t.Errorf("Expected file %q not found in results", name)
			continue
		}
		if gotContent != wantContent {
			t.Errorf("Content mismatch for file %q:\nGot:\n%q\nWant:\n%q", name, gotContent, wantContent)
		}
	}

	// Check if any unexpected files are present in got
	for name := range gotMap {
		if _, ok := wantMap[name]; !ok {
			t.Errorf("Unexpected file %q found in results", name)
		}
	}
}

// --- Test Cases ---

func TestScanEmptyFolder(t *testing.T) {
	tmpDir := createTestDir(t, map[string]string{}) // Create an empty directory
	defer os.RemoveAll(tmpDir)                      // Clean up after the test

	kotlinFiles, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(kotlinFiles.Files) != 0 {
		t.Errorf("Expected 0 Kotlin files in an empty folder, got %d", len(kotlinFiles.Files))
	}
}

func TestScanSingleLevel(t *testing.T) {
	files := map[string]string{
		"MyClass.kt":     "package com.example\nclass MyClass {}",
		"AnotherFile.kt": "fun doSomething() {}",
		"readme.md":      "# Readme",
		"config.json":    "{}",
		"temp.kt.bak":    "backup", // Not a .kt file
	}
	tmpDir := createTestDir(t, files)
	defer os.RemoveAll(tmpDir)

	expectedFiles := []DataFile{
		{Name: "MyClass.kt", Content: []byte("package com.example\nclass MyClass {}")},
		{Name: "AnotherFile.kt", Content: []byte("fun doSomething() {}")},
	}

	kotlinFiles, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Because your `DataFile` struct's `Name` field stores only the base name,
	// and your `ScanFolders` recursive call has an issue (see important note below),
	// this test is designed for a single flat directory scan.
	// For nested directory tests, the current `ScanFolders` implementation would not work as expected.
	compareKotlinFiles(t, kotlinFiles.Files, expectedFiles, tmpDir)
}

func TestScanNestedFolders(t *testing.T) {
	files := map[string]string{
		"src/main/App.kt":               "fun main() {}",
		"src/main/util/Helper.kt":       "class Helper {}",
		"src/test/TestApp.kt":           "import org.junit.Test",
		"src/test/data/TestConfig.json": "{}",        // Non-kotlin in nested
		"config.properties":             "key=value", // Non-kotlin in root
	}
	tmpDir := createTestDir(t, files)
	defer os.RemoveAll(tmpDir)

	// Due to a bug in the provided ScanFolders function,
	// this test case will likely fail to find nested files.
	// The `ScanFolders` recursive call passes the *same* 'folder' argument
	// instead of the path to the sub-directory.
	// Expected files assuming the bug IS fixed and it scans recursively:
	expectedFiles := []DataFile{
		{Name: "App.kt", Content: []byte("fun main() {}")},
		{Name: "Helper.kt", Content: []byte("class Helper {}")},
		{Name: "TestApp.kt", Content: []byte("import org.junit.Test")},
	}

	// If you run this test, it will only find "App.kt" if `ScanFolders`
	// is corrected to recurse into `filepath.Join(folder, file.Name())`.
	// As is, it won't find `Helper.kt` or `TestApp.kt`
	// because `ScanFolders` is called with the original `folder` path for subdirectories.

	kotlinFiles, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// We use the `expectedFiles` with base names here because your `DataFile` stores `Name`, not `Path`.
	// If `DataFile` stored `Path`, the `expectedFiles` would need full relative paths like "src/main/App.kt".
	// The `compareKotlinFiles` helper tries to account for this by normalizing.
	compareKotlinFiles(t, kotlinFiles.Files, expectedFiles, tmpDir)

	// You might want to add an assertion here to explicitly check that
	// `len(kotlinFiles.Files)` is less than `len(expectedFiles)`
	// if you are demonstrating the bug.
}

func TestScanNonExistentFolder(t *testing.T) {
	nonExistentPath := "./non_existent_folder_12345"
	_, err := Scan(nonExistentPath)

	if err == nil {
		t.Error("Expected an error for a non-existent folder, but got none")
	}
	expectedErrMsg := fmt.Sprintf("the folder %s cannot be opened", nonExistentPath)
	if err != nil && !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedErrMsg, err.Error())
	}
}

// TestScanUnreadableFile (Optional - uncomment and modify if you need to test file read errors)
/*
func TestScanUnreadableFile(t *testing.T) {
	// This test is harder to write reliably across different OS and environments
	// as setting unreadable permissions can be tricky and may require root.
	// It's also dependent on how `os.ReadFile` behaves with permissions.
	// The current code immediately returns on the first read error.
	// If you want to test this, you'd create a dummy file and then try to set its permissions
	// to make it unreadable before calling Scan.
	// Example (may not work on all systems):
	// testDir := createTestDir(t, map[string]string{"unreadable.kt": "content"})
	// defer os.RemoveAll(testDir)
	// unreadableFilePath := filepath.Join(testDir, "unreadable.kt")
	//
	// // Attempt to make the file unreadable for the current user (permissions 000)
	// err := os.Chmod(unreadableFilePath, 000)
	// if err != nil {
	// 	t.Skipf("Could not change permissions to test unreadable file: %v", err)
	// }
	//
	// _, err = Scan(testDir)
	// if err == nil {
	// 	t.Error("Expected an error when reading an unreadable file, but got none")
	// }
	// expectedErrMsg := fmt.Sprintf("the file %s cannot be read", "unreadable.kt")
	// if err != nil && !strings.Contains(err.Error(), expectedErrMsg) {
	// 	t.Errorf("Expected error message to contain %q, but got %q", expectedErrMsg, err.Error())
	// }
}
*/
