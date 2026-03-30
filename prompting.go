package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var ErrUserCanceled = errors.New("user canceled")

const rootOption = "(root)"

func promptCategory(ctx context.Context, homeDir string, directories map[string]Category) (string, string, error) {
	keys := make([]string, 0, len(directories))
	for key := range directories {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	selected, err := showRofi(ctx, keys, "Category")
	if err != nil {
		return "", "", err
	}
	if selected == "" {
		return "", "", ErrUserCanceled
	}
	dir, ok := directories[selected]
	if !ok {
		return "", "", fmt.Errorf("invalid category: %s", selected)
	}
	if !dir.PromptSubfolder {
		return selected, "", nil
	}
	subDirs, err := getSubDirectories(filepath.Join(homeDir, dir.Path))
	if err != nil {
		return "", "", err
	}
	if len(subDirs) == 0 {
		return selected, "", nil
	}
	sort.Strings(subDirs)
	options := make([]string, 0, len(subDirs)+1)
	options = append(options, rootOption)
	options = append(options, subDirs...)

	subCat, err := showRofi(ctx, options, selected+" -> ")
	if err != nil {
		return "", "", err
	}

	if subCat == rootOption {
		return selected, "", nil
	}

	return selected, subCat, nil
}

func showRofi(ctx context.Context, options []string, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "rofi", "-dmenu", "-p", prompt)
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() == 1 {
				return "", ErrUserCanceled
			}
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
