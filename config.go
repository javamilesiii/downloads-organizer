package main

import (
	"errors"
	"fmt"
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

func loadConfig(homeDir string) (*Config, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(homeDir, ".config")
	}
	configPath := filepath.Join(configDir, "downloads-organizer", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config from %s: %w", configPath, err)
	}

	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if conf.DownloadsPath == "" {
		conf.DownloadsPath = "Downloads"
	}
	if len(conf.Categories) == 0 {
		return nil, errors.New("no categories configured")
	}
	for name, cat := range conf.Categories {
		if cat.Path == "" {
			return nil, fmt.Errorf("category %q has no path configured", name)
		}
	}
	return conf, nil
}
