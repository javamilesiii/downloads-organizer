package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type Organization struct {
	sourcePath  string
	category    string
	subcategory string
	homeDir     string
	directories map[string]Category
}

func organizeFile(org Organization) error {
	targetDir, exists := org.directories[org.category]
	if !exists {
		return fmt.Errorf("category not found: %s", org.category)
	}
	targetPath := filepath.Join(targetDir.Path, org.subcategory)
	destinationPath := resolveDestination(filepath.Join(org.homeDir, targetPath), filepath.Base(org.sourcePath))

	if err := os.Rename(org.sourcePath, destinationPath); err != nil {
		var linkErr *os.LinkError
		if errors.As(err, &linkErr) && errors.Is(linkErr.Err, syscall.EXDEV) {
			if err := copyFile(org.sourcePath, destinationPath); err != nil {
				return fmt.Errorf("failed to copy file to %s: %w", destinationPath, err)
			}
			if err := os.Remove(org.sourcePath); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", org.sourcePath, err)
			}
			return nil
		}
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}

func isTempFile(filename string, ignoreExts []string) bool {
	if strings.HasPrefix(filepath.Base(filename), ".") {
		return true
	}
	for _, ext := range ignoreExts {
		if !strings.HasPrefix(ext, ".") {
			if strings.HasPrefix(filepath.Base(filename), ext) {
				return true
			}
		} else if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func getSubDirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", path, err)
	}
	var subDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry.Name())
		}
	}
	return subDirs, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(sourcePath string, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error occurred while opening file %s: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("error occurred while creating file %s: %w", destinationPath, err)
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		os.Remove(destinationPath)
		return fmt.Errorf("error occurred while copying file: %w", err)
	}
	if err := destinationFile.Close(); err != nil {
		os.Remove(destinationPath)
		return fmt.Errorf("error occurred while closing destinationFile %s: %w", destinationPath, err)
	}
	return nil
}

func resolveDestination(dir, filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	path := filepath.Join(dir, filename)
	if !pathExists(path) {
		return path
	}
	for i := 1; ; i++ {
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, i, ext))
		if !pathExists(newPath) {
			return newPath
		}
	}
}
