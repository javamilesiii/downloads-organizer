package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var homeDir = getHomeDir()
var directories = map[string]string{
	"School":   "Documents/School",
	"Coding":   "Development/Downloads",
	"Scouts":   "Documents/Scouts",
	"Music":    "Music",
	"Videos":   "Videos",
	"Pictures": "Pictures",
}

func main() {
	if err := setup(); err != nil {
		log.Println("Error while setup: ", err)
		return
	}
	if err := watchDownloads(); err != nil {
		log.Println("Error while watching downloads: ", err)
	}
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get home directory: %v", err)
	}
	return homeDir
}

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
	destinationPath := filepath.Join(homeDir, targetDir, filepath.Base(sourcePath))
	if err := os.Rename(sourcePath, destinationPath); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}

func setup() error {
	for _, value := range directories {
		path := filepath.Join(homeDir, value)
		if _, err := os.Stat(path); err == nil {
			log.Println("Directory already exists: ", path)
			continue
		} else if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		log.Println("Directory successfully created: ", path)
	}
	downloadDir := filepath.Join(homeDir, "Downloads")
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(downloadDir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func promptCategory() (string, error) {
	cmd := exec.Command("rofi", "-dmenu", "-p", "Category:", "-theme-str", "window {width: 30%;}")
	keys := make([]string, 0, len(directories))
	for key := range directories {
		keys = append(keys, key)
	}
	options := strings.Join(keys, "\n")
	cmd.Stdin = strings.NewReader(options)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	selected := strings.TrimSpace(string(output))
	if selected == "" {
		return "", fmt.Errorf("no category selected")
	}
	if _, ok := directories[selected]; !ok {
		return "", fmt.Errorf("invalid category: %s", selected)
	}
	return selected, nil
}

func watchDownloads() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	downloadsPath := filepath.Join(homeDir, "Downloads")
	if err := watcher.Add(downloadsPath); err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				if isTempFile(event.Name) {
					continue
				}
				category, err := promptCategory()
				if err != nil {
					log.Println("Error while getting category: ", err)
					continue
				}

				if err := organizeFile(event.Name, category); err != nil {
					log.Println("Error while moving file: ", err)
					continue
				}
			}
		case err := <-watcher.Errors:
			log.Println("Error:", err)
		}
	}
}

func isTempFile(filename string) bool {
	tempExtensions := []string{
		".crdownload",
		".part",
		".tmp",
		".opdownload",
		".download",
		".temp",
	}
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
