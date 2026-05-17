package audio

import (
	"strings"
	"testing"
)

func TestParsePlaylistPLS(t *testing.T) {
	playlist := `[playlist]
NumberOfEntries=2
File1=https://example.com/first.mp3
Title1=First
File2=https://example.com/second.mp3
`
	got := parsePlaylist(strings.NewReader(playlist))
	if got != "https://example.com/first.mp3" {
		t.Fatalf("got %q", got)
	}
}

func TestParsePlaylistM3U(t *testing.T) {
	playlist := `#EXTM3U
#EXTINF:-1,Station
https://example.com/live.aac
`
	got := parsePlaylist(strings.NewReader(playlist))
	if got != "https://example.com/live.aac" {
		t.Fatalf("got %q", got)
	}
}

func TestParsePlaylistIgnoresNoise(t *testing.T) {
	playlist := `

# comment
[playlist]
Title1=No URL Here
not a url
https://example.com/stream.ogg
`
	got := parsePlaylist(strings.NewReader(playlist))
	if got != "https://example.com/stream.ogg" {
		t.Fatalf("got %q", got)
	}
}
