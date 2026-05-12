package directory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type somaPlaylist struct {
	URL     string `json:"url"`
	Format  string `json:"format"`
	Quality string `json:"quality"`
}

func (c Client) SomaFM(ctx context.Context) ([]Station, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.somafm.com/channels.json", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("somafm: %s", resp.Status)
	}

	var payload struct {
		Channels []struct {
			Title       string         `json:"title"`
			Description string         `json:"description"`
			Genre       string         `json:"genre"`
			Listeners   int            `json:"listeners,string"`
			LastPlaying string         `json:"lastPlaying"`
			Playlists   []somaPlaylist `json:"playlists"`
		} `json:"channels"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	out := make([]Station, 0, len(payload.Channels))
	for _, ch := range payload.Channels {
		stream := firstPlaylist(ch.Playlists)
		if stream == "" {
			continue
		}
		out = append(out, Station{
			Name:        ch.Title,
			Description: strings.TrimSpace(ch.Description),
			Genre:       ch.Genre,
			URL:         stream,
			Source:      "SomaFM",
			Listeners:   ch.Listeners,
			NowPlaying:  ch.LastPlaying,
		})
	}
	return out, nil
}

func firstPlaylist(playlists []somaPlaylist) string {
	for _, pref := range []string{"mp3", "aac"} {
		for _, p := range playlists {
			if strings.EqualFold(p.Format, pref) && p.URL != "" {
				return p.URL
			}
		}
	}
	if len(playlists) > 0 {
		return playlists[0].URL
	}
	return ""
}
