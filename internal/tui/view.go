package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) View() string {
	var b strings.Builder
	contentWidth := max(40, m.width-4)
	b.WriteString(gradientLogo(logo))
	b.WriteString("\n\n")
	b.WriteString(m.styles.header.Render("[1] presets  [2] SomaFM  [3] Xiph  [4] my stations  [/] tune  [ctrl+p] preset  [enter] play/add  [s] stop  [q] quit"))
	b.WriteString("\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(m.styles.panel.Width(contentWidth).Render(m.visualizer.View(m.styles)))
	b.WriteString("\n")
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.statusLine())
	if m.err != "" {
		b.WriteString("\n" + m.styles.error.Render(m.err))
	}
	return m.styles.frame.Render(b.String())
}

func (m Model) statusLine() string {
	if m.playing != nil {
		return m.styles.status.Render("now: " + m.playing.Name + " <" + m.playing.URL + ">")
	}
	return m.styles.status.Render(m.status)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
