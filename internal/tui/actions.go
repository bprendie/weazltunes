package tui

import (
	"context"
	"fmt"

	"github.com/bprendie/weazltunes/internal/audio"
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

func (m *Model) startRenameSelected() {
	item, ok := m.list.SelectedItem().(stationItem)
	if !ok || (m.mode != modePresets && m.mode != modeMyStations) {
		m.err = "select a preset or my station to rename"
		return
	}
	st := directoryStation(item)
	m.renaming = &st
	m.input.SetValue(st.Name)
	m.input.Prompt = "rename > "
	m.input.Focus()
	m.status = "enter a new station name"
	m.err = ""
}

func (m *Model) finishRename() {
	if m.renaming == nil {
		return
	}
	name := m.input.Value()
	renamed := false
	if m.mode == modePresets {
		renamed = m.cfg.RenamePreset(m.renaming.URL, name)
	} else if m.mode == modeMyStations {
		renamed = m.cfg.RenameMyStation(m.renaming.URL, name)
	}
	if !renamed {
		m.err = "rename failed"
		return
	}
	if err := config.Save(m.cfg); err != nil {
		m.err = err.Error()
		return
	}
	m.input.SetValue("")
	m.input.Prompt = tunePrompt
	m.input.Blur()
	m.renaming = nil
	if m.mode == modePresets {
		m.showPresets()
	} else {
		m.showMyStations()
	}
	m.status = "renamed"
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
