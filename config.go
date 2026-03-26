package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DownloadsPath string              `yaml:"downloads_path"`
	Categories    map[string]Category `yaml:"categories"`
	ignoreFiles   []string            `yaml:"ignore_files"`
}

type Category struct {
	path            string `yaml:"path"`
	promptSubfolder bool   `yaml:"promptSubfolder"`
}

var config Config

func loadConfig() {
	data, err := os.ReadFile(filepath.Join(homeDir, ".config/downloads-organizer/config.yaml"))
	if err != nil {
		panic(err)
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}
	config = conf
}

func reloadConfig() {
	loadConfig()
}

func getConfig() Config {
	return config
}

func getCategories() map[string]Category {
	return config.Categories
}

func getDownloadsPath() string {
	return config.DownloadsPath
}

func getIgnoreFiles() []string {
	return config.ignoreFiles
}
