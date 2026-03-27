package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func promptCategory() (string, string, error) {
	keys := make([]string, 0, len(directories))
	for key := range directories {
		keys = append(keys, key)
	}

	selected, err := showRofi(keys, "Category")
	if err != nil {
		return "", "", err
	}
	if selected == "" {
		return "", "", fmt.Errorf("no category selected")
	}
	if _, ok := directories[selected]; !ok {
		return "", "", fmt.Errorf("invalid category: %s", selected)
	}
	if !directories[selected].PromptSubfolder {
		return selected, "", nil
	}
	subDirs, err := getSubDirectories(filepath.Join(homeDir, selected))
	if err != nil || len(subDirs) == 0 {
		return selected, "", nil
	}

	options := append([]string{"(root)"}, subDirs...)

	subCat, err := showRofi(options, selected+" -> ")
	if err != nil {
		return "", "", err
	}

	if subCat == "(root)" {
		return selected, "", nil
	}

	return selected, subCat, nil
}

func showRofi(options []string, prompt string) (string, error) {
	cmd := exec.Command("rofi", "-dmenu", "-p", "Category")
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
