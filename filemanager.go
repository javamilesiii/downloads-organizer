package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func organizeFile(sourcePath, category string, subcategory string) error {
	targetDir, exists := directories[category]
	if !exists {
		return fmt.Errorf("category not found: %s", category)
	}
	if !dirExists(filepath.Join(targetDir.Path, subcategory)) && subcategory != "" {
		return fmt.Errorf("subcategory not found")
	}
	targetDir.Path = filepath.Join(targetDir.Path, subcategory)
	if !dirExists(sourcePath) {
		return fmt.Errorf("file not found: %s", sourcePath)
	}
	destinationPath := filepath.Join(homeDir, targetDir.Path, filepath.Base(sourcePath))
	if err := os.Rename(sourcePath, destinationPath); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}

func isTempFile(filename string) bool {
	tempExtensions := config.getIgnoreFiles()
	if strings.HasPrefix(filepath.Base(filename), ".") {
		return true
	}
	if strings.Contains(filename, "Unconfirmed") {
		return true
	}
	for _, ext := range tempExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func getSubDirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var subDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry.Name())
		}
	}
	return subDirs, nil
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}
