package tui

import (
	"net/url"
	"path"
	"strings"

	"github.com/bprendie/weazltunes/internal/config"
	"github.com/bprendie/weazltunes/internal/directory"
)

func presetStations(cfg config.Config) []directory.Station {
	return configStations(cfg.Presets, "Preset")
}

func myStations(cfg config.Config) []directory.Station {
	return configStations(cfg.MyStations, "My Station")
}

func configStations(items []config.Preset, source string) []directory.Station {
	out := make([]directory.Station, 0, len(items))
	for _, p := range items {
		out = append(out, directory.Station{Name: p.Name, URL: p.URL, Source: source, Description: p.URL})
	}
	return out
}

func stationFromURL(raw string) (directory.Station, bool) {
	raw = strings.TrimSpace(raw)
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return directory.Station{}, false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return directory.Station{}, false
	}
	name := stationName(parsed)
	return directory.Station{Name: name, URL: raw, Source: "My Station", Description: raw}, true
}

func stationName(u *url.URL) string {
	if strings.Contains(u.Host, "youtube.com") && u.Query().Get("v") != "" {
		return "youtube " + u.Query().Get("v")
	}
	base := path.Base(u.Path)
	if base == "." || base == "/" || base == "" {
		return u.Host
	}
	if base == "watch" {
		return u.Host
	}
	return strings.TrimSuffix(base, path.Ext(base))
}

func presetFromStation(st directory.Station) config.Preset {
	return config.Preset{Name: st.Name, URL: st.URL}
}
