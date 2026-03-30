package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	homeDir := getHomeDir()
	config, err := loadConfig(homeDir)
	if err != nil {
		log.Fatalf("error while loading config: %v", err)
	}
	downloadsPath, err := setup(homeDir, config)
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}
	if err := watchDownloads(homeDir, downloadsPath, config); err != nil {
		log.Fatalf("error while watching downloads: %v", err)
	}
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get home directory: %v", err)
	}
	return homeDir
}

func setup(homeDir string, config *Config) (string, error) {
	for _, dir := range config.Categories {
		path := filepath.Join(homeDir, dir.Path)
		if _, err := os.Stat(path); err == nil {
			continue
		} else if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return "", fmt.Errorf("creating directory %s: %w", path, err)
			}
			log.Printf("directory successfully created: %s", path)
		} else if err != nil {
			return "", fmt.Errorf("checking directory %s: %w", path, err)
		}
	}
	downloadsPath := filepath.Join(homeDir, config.DownloadsPath)
	if _, err := os.Stat(downloadsPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(downloadsPath, 0755); err != nil {
			return "", fmt.Errorf("creating downloads directory %s: %w", downloadsPath, err)
		}
	} else if err != nil {
		return "", fmt.Errorf("checking downloads directory %s: %w", downloadsPath, err)
	}
	return downloadsPath, nil
}

func watchDownloads(homeDir string, downloadsPath string, config *Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)
	defer watcher.Close()

	if err := watcher.Add(downloadsPath); err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			// Files are handled sequentially: rofi blocks the loop, so concurrent
			// downloads will que up and prompt one after another.
			if event.Has(fsnotify.Create) {
				if isTempFile(event.Name, config.IgnoreFiles) {
					continue
				}
				category, subcategory, err := promptCategory(ctx, homeDir, config.Categories)
				if err != nil {
					if errors.Is(err, ErrUserCanceled) {
						continue
					}
					log.Printf("error while getting category: %v", err)
					continue
				}
				org := Organization{
					sourcePath:  event.Name,
					category:    category,
					subcategory: subcategory,
					homeDir:     homeDir,
					directories: config.Categories,
				}
				if err := organizeFile(org); err != nil {
					log.Printf("error while moving file: %v", err)
					continue
				}
			}
		case err := <-watcher.Errors:
			log.Printf("error: %v", err)
		case <-sigChan:
			cancel()
			log.Println("shutting down...")
			return nil
		}
	}
}
