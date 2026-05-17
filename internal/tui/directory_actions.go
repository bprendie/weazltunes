package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

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
