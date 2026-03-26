package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var homeDir = getHomeDir()
var directories map[string]Category
var downloadsPath string

var config Config

func main() {
	loadConfig()
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

func setup() error {
	directories = config.getCategories()
	downloadsPath = config.getDownloadsPath()
	for _, value := range directories {
		path := filepath.Join(homeDir, value.Path)
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
	if _, err := os.Stat(downloadsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(downloadsPath, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func watchDownloads() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

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
				if _, err := os.Stat(event.Name); err != nil {
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
