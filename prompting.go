package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func promptCategory() (string, error) {
	cmd := exec.Command("rofi", "-dmenu", "-p", "Category")
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
