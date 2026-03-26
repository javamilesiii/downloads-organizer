package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DownloadsPath string              `yaml:"download_dir"`
	Categories    map[string]Category `yaml:"categories"`
	IgnoreFiles   []string            `yaml:"ignore_files"`
}

type Category struct {
	Path            string `yaml:"path"`
	PromptSubfolder bool   `yaml:"prompt_subfolder"`
}

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

func (c Config) getCategories() map[string]Category {
	return c.Categories
}

func (c Config) getDownloadsPath() string {
	return c.DownloadsPath
}

func (c Config) getIgnoreFiles() []string {
	return c.IgnoreFiles
}
