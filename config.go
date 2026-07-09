package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	Lang string `json:"lang,omitempty"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "diskus", "config.json"), nil
}

func loadConfig() config {
	var c config
	p, err := configPath()
	if err != nil {
		return c
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return c
	}
	_ = json.Unmarshal(data, &c)
	return c
}

func saveConfig(c config) (string, error) {
	p, err := configPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(p, data, 0o644); err != nil {
		return "", err
	}
	return p, nil
}
