package audio

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

func ResolveStreamURL(ctx context.Context, raw string) string {
	if !looksLikePlaylist(raw) {
		return raw
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return raw
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return raw
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return raw
	}
	if stream := parsePlaylist(resp.Body); stream != "" {
		return stream
	}
	return raw
}

func looksLikePlaylist(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	ext := strings.ToLower(path.Ext(u.Path))
	return ext == ".pls" || ext == ".m3u" || ext == ".m3u8"
}

func parsePlaylist(r io.Reader) string {
	scanner := bufio.NewScanner(io.LimitReader(r, 1<<20))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}
		if key, value, ok := strings.Cut(line, "="); ok {
			if strings.HasPrefix(strings.ToLower(key), "file") && isHTTP(value) {
				return strings.TrimSpace(value)
			}
			continue
		}
		if isHTTP(line) {
			return line
		}
	}
	return ""
}

func isHTTP(raw string) bool {
	return strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://")
}
