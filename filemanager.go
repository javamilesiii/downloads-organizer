package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func organizeFile(sourcePath, category string) error {
	targetDir, exists := directories[category]
	if !exists {
		return fmt.Errorf("category not found: %s", category)
	}
	_, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", sourcePath)
	} else if err != nil {
		return err
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
