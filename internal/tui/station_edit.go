package tui

import (
	"github.com/bprendie/weazltunes/internal/config"
	"github.com/bprendie/weazltunes/internal/directory"
)

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
	renamed := false
	if m.mode == modePresets {
		renamed = m.cfg.RenamePreset(m.renaming.URL, m.input.Value())
	} else if m.mode == modeMyStations {
		renamed = m.cfg.RenameMyStation(m.renaming.URL, m.input.Value())
	}
	if !renamed || !m.saveConfig() {
		return
	}
	m.input.SetValue("")
	m.input.Prompt = tunePrompt
	m.input.Blur()
	m.renaming = nil
	m.refreshEditableList()
	m.status = "renamed"
}

func (m *Model) deleteSelected() {
	item, ok := m.list.SelectedItem().(stationItem)
	if !ok || (m.mode != modePresets && m.mode != modeMyStations) {
		m.err = "select a preset or my station to delete"
		return
	}
	st := directoryStation(item)
	deleted := false
	if m.mode == modePresets {
		deleted = m.cfg.DeletePreset(st.URL)
	} else {
		deleted = m.cfg.DeleteMyStation(st.URL)
	}
	if !deleted || !m.saveConfig() {
		return
	}
	if m.playing != nil && m.playing.URL == st.URL {
		m.stop()
	}
	m.refreshEditableList()
	m.status = "deleted " + st.Name
}

func (m *Model) saveConfig() bool {
	if err := config.Save(m.cfg); err != nil {
		m.err = err.Error()
		return false
	}
	return true
}

func (m *Model) refreshEditableList() {
	if m.mode == modePresets {
		m.showPresets()
	} else {
		m.showMyStations()
	}
}
