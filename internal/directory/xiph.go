package directory

import (
	"context"
	"encoding/xml"
	"html"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

func (c Client) XiphGenres(ctx context.Context) ([]string, error) {
	body, err := c.getString(ctx, "https://dir.xiph.org/genres")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(?i)<a[^>]+href="/genres/([^"#?]+)"[^>]*>([^<]+)</a>`)
	seen := map[string]bool{}
	var genres []string
	for _, m := range re.FindAllStringSubmatch(body, -1) {
		g := html.UnescapeString(strings.TrimSpace(m[2]))
		if g == "" || seen[strings.ToLower(g)] {
			continue
		}
		seen[strings.ToLower(g)] = true
		genres = append(genres, g)
	}
	sort.Strings(genres)
	return genres, nil
}

func (c Client) XiphSearch(ctx context.Context, query string) ([]Station, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return c.xiphXML(ctx, "")
	}
	if stations, err := c.xiphGenrePage(ctx, query); err == nil && len(stations) > 0 {
		return stations, nil
	}
	return c.xiphXML(ctx, query)
}

func (c Client) xiphXML(ctx context.Context, query string) ([]Station, error) {
	body, err := c.getString(ctx, "https://dir.xiph.org/yp.xml")
	if err != nil {
		return nil, err
	}
	var payload struct {
		Entries []struct {
			ServerName        string `xml:"server_name"`
			ListenURL         string `xml:"listen_url"`
			Genre             string `xml:"genre"`
			CurrentSong       string `xml:"current_song"`
			ServerDescription string `xml:"server_description"`
			Listeners         int    `xml:"listeners"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal([]byte(body), &payload); err != nil {
		return nil, err
	}
	return filterXiphEntries(payload.Entries, query), nil
}

func (c Client) xiphGenrePage(ctx context.Context, genre string) ([]Station, error) {
	body, err := c.getString(ctx, "https://dir.xiph.org/genres/"+url.PathEscape(genre))
	if err != nil {
		return nil, err
	}
	blocks := regexp.MustCompile(`(?is)<h5[^>]*>(.*?)</h5>(.*?)(?:<h5|$)`).FindAllStringSubmatch(body, -1)
	var out []Station
	for _, block := range blocks {
		name := strings.TrimSpace(stripTags(block[1]))
		link := firstHref(block[2])
		if name == "" || link == "" {
			continue
		}
		out = append(out, Station{
			Name:        html.UnescapeString(name),
			Description: strings.TrimSpace(html.UnescapeString(stripTags(block[2]))),
			Genre:       genre,
			URL:         link,
			Source:      "Xiph",
		})
	}
	return out, nil
}

func filterXiphEntries(entries []struct {
	ServerName        string `xml:"server_name"`
	ListenURL         string `xml:"listen_url"`
	Genre             string `xml:"genre"`
	CurrentSong       string `xml:"current_song"`
	ServerDescription string `xml:"server_description"`
	Listeners         int    `xml:"listeners"`
}, query string) []Station {
	q := strings.ToLower(query)
	var out []Station
	for _, e := range entries {
		haystack := strings.ToLower(e.ServerName + " " + e.Genre + " " + e.ServerDescription + " " + e.CurrentSong)
		if e.ListenURL == "" || (q != "" && !strings.Contains(haystack, q)) {
			continue
		}
		out = append(out, Station{Name: e.ServerName, Description: e.ServerDescription, Genre: e.Genre, URL: e.ListenURL, Source: "Xiph", Listeners: e.Listeners, NowPlaying: e.CurrentSong})
		if len(out) >= 200 {
			break
		}
	}
	return out
}
