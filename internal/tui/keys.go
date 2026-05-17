package tui

import tea "github.com/charmbracelet/bubbletea"

func handleKey(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.stopMeter()
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
		promoteSelectedOrInput(&m)
	case "ctrl+r":
		m.startRenameSelected()
	case "ctrl+d":
		m.deleteSelected()
	case "[":
		m.moveSelected(-1)
	case "]":
		m.moveSelected(1)
	case " ":
		m.togglePause()
	case "enter":
		return handleEnter(m)
	case "/":
		m.input.Focus()
	case "esc":
		m.renaming = nil
		m.input.SetValue("")
		m.input.Prompt = tunePrompt
		m.input.Blur()
		return m, noop
	case "s":
		m.stop()
	}
	return m, nil
}

func handleEnter(m Model) (Model, tea.Cmd) {
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

func promoteSelectedOrInput(m *Model) {
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

func noop() tea.Msg {
	return nil
}
