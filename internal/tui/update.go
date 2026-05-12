package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		next, cmd := m.handleKey(msg)
		m = next
		if cmd != nil {
			return m, cmd
		}
	case loadedMsg:
		m.setStations(msg.stations)
		m.status = msg.status
		m.err = ""
		m.input.Blur()
	case errMsg:
		m.err = msg.err.Error()
	case tickMsg:
		m.visualizer.Step(m.playing != nil && !m.paused)
		cmds = append(cmds, tick())
	}

	next, cmd := m.updateFocused(msg)
	m = next
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.player.Stop()
		return m, tea.Quit
	case "1":
		m.showPresets()
	case "2":
		m.mode = modeSoma
		m.status = "loading SomaFM"
		m.err = ""
		return m, m.loadSoma()
	case "3":
		m.mode = modeXiph
		m.status = "enter a Xiph genre/search or stream URL"
		m.err = ""
		m.input.Focus()
	case "4":
		m.showMyStations()
	case "ctrl+p":
		m.promoteSelectedOrInput()
	case "ctrl+r":
		m.startRenameSelected()
	case "v":
		m.status = "visualizer: " + m.visualizer.Toggle()
	case " ":
		m.togglePause()
	case "enter":
		return m.handleEnter()
	case "/":
		m.input.Focus()
	case "esc":
		m.renaming = nil
		m.input.SetValue("")
		m.input.Prompt = tunePrompt
		m.input.Blur()
	case "s":
		m.stop()
	}
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	if m.input.Focused() {
		if m.renaming != nil {
			m.finishRename()
			return m, nil
		}
		if st, ok := stationFromURL(m.input.Value()); ok {
			m.saveMyStation(st)
			m.showMyStations()
			m.play(st)
			m.input.SetValue("")
			m.input.Blur()
			return m, nil
		}
		m.status = "searching Xiph"
		m.err = ""
		return m, m.searchXiph(m.input.Value())
	}
	if item, ok := m.list.SelectedItem().(stationItem); ok {
		m.play(directoryStation(item))
	}
	return m, nil
}

func (m *Model) promoteSelectedOrInput() {
	if st, ok := stationFromURL(m.input.Value()); ok {
		m.saveMyStation(st)
		m.promotePreset(st)
		m.input.SetValue("")
		m.input.Blur()
		m.showPresets()
		m.status = "added to presets"
		return
	}
	if item, ok := m.list.SelectedItem().(stationItem); ok {
		m.promotePreset(directoryStation(item))
		m.showPresets()
		m.status = "added to presets"
		return
	}
	if m.playing != nil {
		m.promotePreset(*m.playing)
		m.showPresets()
		m.status = "added to presets"
	}
}

func (m *Model) resize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(max(20, width-8), max(8, height-18))
	m.input.Width = max(20, width-20)
}

func (m Model) updateFocused(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == " " {
		return m, nil
	}
	if m.input.Focused() {
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
