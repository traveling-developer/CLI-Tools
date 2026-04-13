package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFilePath = ".config/unsplash/config.json"

type apiConfig struct {
	AccessKey string `json:"access_key"`
}

func resolveAPIKey(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	if key := os.Getenv("UNSPLASH_ACCESS_KEY"); key != "" {
		return key, nil
	}

	if key, err := keyFromFile(); err == nil {
		return key, nil
	}

	return "", fmt.Errorf(
		"no API key found — provide one via:\n" +
			"  --api-key <key>\n" +
			"  UNSPLASH_ACCESS_KEY env var\n" +
			"  ~/.config/unsplash/config.json  →  { \"access_key\": \"<key>\" }",
	)
}

func keyFromFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	f, err := os.Open(filepath.Join(home, configFilePath))
	if err != nil {
		return "", err
	}
	defer f.Close()

	var cfg apiConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return "", err
	}

	if cfg.AccessKey == "" {
		return "", fmt.Errorf("access_key is empty in config file")
	}

	return cfg.AccessKey, nil
}
