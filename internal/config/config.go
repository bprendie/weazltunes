package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Presets    []Preset `json:"presets"`
	MyStations []Preset `json:"my_stations"`
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

func (c *Config) SaveMyStation(st Preset) {
	st.Name = strings.TrimSpace(st.Name)
	st.URL = strings.TrimSpace(st.URL)
	if st.URL == "" {
		return
	}
	c.MyStations = upsert(c.MyStations, st)
}

func (c *Config) PromotePreset(st Preset) {
	st.Name = strings.TrimSpace(st.Name)
	st.URL = strings.TrimSpace(st.URL)
	if st.URL == "" {
		return
	}
	presets := upsert(c.Presets, st)
	c.Presets = presets
	if len(c.Presets) > 5 {
		c.Presets = c.Presets[:5]
	}
}

func (c *Config) RenamePreset(url, name string) bool {
	return renameByURL(c.Presets, url, name)
}

func (c *Config) RenameMyStation(url, name string) bool {
	return renameByURL(c.MyStations, url, name)
}

func (c *Config) DeletePreset(url string) bool {
	next, deleted := deleteByURL(c.Presets, url)
	c.Presets = next
	return deleted
}

func (c *Config) DeleteMyStation(url string) bool {
	next, deleted := deleteByURL(c.MyStations, url)
	c.MyStations = next
	return deleted
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

func upsert(stations []Preset, st Preset) []Preset {
	out := []Preset{st}
	for _, existing := range stations {
		if strings.EqualFold(existing.URL, st.URL) {
			continue
		}
		out = append(out, existing)
	}
	return out
}

func renameByURL(stations []Preset, url, name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	for i := range stations {
		if strings.EqualFold(stations[i].URL, url) {
			stations[i].Name = name
			return true
		}
	}
	return false
}

func deleteByURL(stations []Preset, url string) ([]Preset, bool) {
	out := stations[:0]
	deleted := false
	for _, station := range stations {
		if strings.EqualFold(station.URL, url) {
			deleted = true
			continue
		}
		out = append(out, station)
	}
	return out, deleted
}
