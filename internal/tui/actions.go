package tui

import (
	"context"
	"fmt"

	"github.com/bprendie/weazltunes/internal/audio"
	"github.com/bprendie/weazltunes/internal/directory"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) showPresets() {
	m.mode = modePresets
	m.list.Title = "top 8 presets"
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
	m.paused = false
	m.startMeter(st.URL)
	m.status = "playing " + st.Name
	m.err = ""
}

func (m *Model) togglePause() {
	paused, err := m.player.TogglePause()
	if err != nil {
		m.err = err.Error()
		return
	}
	m.paused = paused
	if paused {
		m.stopMeter()
		m.status = "paused"
		return
	}
	if m.playing != nil {
		m.startMeter(m.playing.URL)
		m.status = "playing " + m.playing.Name
	}
}

func (m *Model) stop() {
	m.player.Stop()
	m.stopMeter()
	m.playing = nil
	m.paused = false
	m.status = "stopped"
}

func (m *Model) startMeter(url string) {
	m.stopMeter()
	meter, err := audio.StartMeter(url)
	if err != nil {
		m.energy = audio.Sample{}
		return
	}
	m.meter = meter
}

func (m *Model) stopMeter() {
	if m.meter != nil {
		m.meter.Stop()
	}
	m.meter = nil
	m.energy = audio.Sample{}
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
