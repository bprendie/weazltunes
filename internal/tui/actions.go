package tui

import (
	"context"
	"fmt"

	"github.com/bprendie/weazltunes/internal/config"
	"github.com/bprendie/weazltunes/internal/directory"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) showPresets() {
	m.mode = modePresets
	m.list.Title = "top 5 presets"
	m.setStations(presetStations(m.cfg))
	m.status = "presets"
	m.err = ""
}

func (m *Model) showMyStations() {
	m.mode = modeMyStations
	m.list.Title = "my stations"
	m.setStations(myStations(m.cfg))
	m.status = "my stations"
	m.err = ""
}

func (m *Model) play(st directory.Station) {
	if err := m.player.Play(st.URL); err != nil {
		m.err = err.Error()
		return
	}
	m.playing = &st
	m.status = "playing " + st.Name
	m.err = ""
}

func (m *Model) stop() {
	m.player.Stop()
	m.playing = nil
	m.status = "stopped"
}

func (m *Model) saveMyStation(st directory.Station) {
	m.cfg.SaveMyStation(presetFromStation(st))
	if err := config.Save(m.cfg); err != nil {
		m.err = err.Error()
		return
	}
	m.status = "saved to my stations"
}

func (m *Model) promotePreset(st directory.Station) {
	m.cfg.PromotePreset(presetFromStation(st))
	if err := config.Save(m.cfg); err != nil {
		m.err = err.Error()
		return
	}
	m.status = "added to presets"
}

func (m *Model) setStations(stations []directory.Station) {
	items := make([]list.Item, 0, len(stations))
	for _, st := range stations {
		items = append(items, stationItem{station: st})
	}
	m.list.SetItems(items)
}

func (m Model) loadSoma() tea.Cmd {
	return func() tea.Msg {
		stations, err := m.client.SomaFM(context.Background())
		if err != nil {
			return errMsg{err}
		}
		return loadedMsg{stations: stations, status: fmt.Sprintf("loaded %d SomaFM channels", len(stations))}
	}
}

func (m Model) searchXiph(q string) tea.Cmd {
	return func() tea.Msg {
		stations, err := m.client.XiphSearch(context.Background(), q)
		if err != nil {
			return errMsg{err}
		}
		return loadedMsg{stations: stations, status: fmt.Sprintf("loaded %d Xiph streams", len(stations))}
	}
}

func directoryStation(item stationItem) directory.Station {
	return item.station
}
