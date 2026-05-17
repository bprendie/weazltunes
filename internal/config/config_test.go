package config

import "testing"

func TestPromotePresetCapsAtEightAndMovesDuplicateToTop(t *testing.T) {
	cfg := Config{}
	for i := 0; i < PresetLimit+2; i++ {
		cfg.PromotePreset(Preset{Name: string(rune('A' + i)), URL: "https://example.com/" + string(rune('a'+i))})
	}
	if len(cfg.Presets) != PresetLimit {
		t.Fatalf("expected %d presets, got %d", PresetLimit, len(cfg.Presets))
	}

	duplicate := cfg.Presets[3]
	cfg.PromotePreset(duplicate)
	if len(cfg.Presets) != PresetLimit {
		t.Fatalf("duplicate changed preset count: %d", len(cfg.Presets))
	}
	if cfg.Presets[0].URL != duplicate.URL {
		t.Fatalf("duplicate was not moved to top: got %q want %q", cfg.Presets[0].URL, duplicate.URL)
	}
}

func TestSaveMyStationUpsertsByURL(t *testing.T) {
	cfg := Config{}
	cfg.SaveMyStation(Preset{Name: "Old", URL: "https://example.com/stream"})
	cfg.SaveMyStation(Preset{Name: "New", URL: "https://example.com/stream"})

	if len(cfg.MyStations) != 1 {
		t.Fatalf("expected one station, got %d", len(cfg.MyStations))
	}
	if cfg.MyStations[0].Name != "New" {
		t.Fatalf("expected latest station name, got %q", cfg.MyStations[0].Name)
	}
}

func TestRenameAndDeleteByURL(t *testing.T) {
	cfg := Config{Presets: []Preset{{Name: "Old", URL: "https://example.com/stream"}}}

	if !cfg.RenamePreset("https://example.com/stream", "New") {
		t.Fatal("expected rename to succeed")
	}
	if cfg.Presets[0].Name != "New" {
		t.Fatalf("rename did not update name: %q", cfg.Presets[0].Name)
	}
	if cfg.RenamePreset("https://example.com/stream", "   ") {
		t.Fatal("expected blank rename to fail")
	}
	if !cfg.DeletePreset("https://example.com/stream") {
		t.Fatal("expected delete to succeed")
	}
	if len(cfg.Presets) != 0 {
		t.Fatalf("expected no presets after delete, got %d", len(cfg.Presets))
	}
}

func TestMoveByURLBounds(t *testing.T) {
	cfg := Config{MyStations: []Preset{
		{Name: "One", URL: "https://example.com/1"},
		{Name: "Two", URL: "https://example.com/2"},
		{Name: "Three", URL: "https://example.com/3"},
	}}

	if next, ok := cfg.MoveMyStation("https://example.com/2", -1); !ok || next != 0 {
		t.Fatalf("expected move up to index 0, got index=%d ok=%v", next, ok)
	}
	if cfg.MyStations[0].Name != "Two" {
		t.Fatalf("expected Two at top, got %q", cfg.MyStations[0].Name)
	}
	if _, ok := cfg.MoveMyStation("https://example.com/2", -1); ok {
		t.Fatal("expected moving beyond top to fail")
	}
	if _, ok := cfg.MoveMyStation("https://example.com/missing", 1); ok {
		t.Fatal("expected missing URL move to fail")
	}
}
