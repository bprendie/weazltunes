package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Presets []Preset `json:"presets"`
}

type Preset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func Load() (Config, error) {
	cfg := defaultConfig()
	path, err := Path()
	if err != nil {
		return cfg, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return cfg, err
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, Save(cfg)
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	if len(cfg.Presets) == 0 {
		cfg.Presets = defaultConfig().Presets
	}
	return cfg, nil
}

func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o600)
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "weazltunes", "config.json"), nil
}

func defaultConfig() Config {
	return Config{
		Presets: []Preset{
			{Name: "Groove Salad", URL: "https://somafm.com/groovesalad.pls"},
			{Name: "Drone Zone", URL: "https://somafm.com/dronezone.pls"},
			{Name: "DEF CON Radio", URL: "https://somafm.com/defcon.pls"},
			{Name: "Indie Pop Rocks", URL: "https://somafm.com/indiepop.pls"},
			{Name: "Deep Space One", URL: "https://somafm.com/deepspaceone.pls"},
		},
	}
}
