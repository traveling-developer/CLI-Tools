package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultFeedsPath = ".config/rss/feeds"

func resolveFeedsFile(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, defaultFeedsPath)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return "", fmt.Errorf(
		"no feeds file found — provide one via:\n" +
			"  --feeds <file>\n" +
			"  ~/.config/rss/feeds  (one URL per line)",
	)
}

func readFeedURLs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening feeds file: %w", err)
	}
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading feeds file: %w", err)
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("feeds file %q contains no URLs", path)
	}
	return urls, nil
}
