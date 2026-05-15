package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func (m Model) View() string {
	var b strings.Builder
	contentWidth := max(40, m.width-4)
	if contentWidth >= maxLineWidth(logo) {
		b.WriteString(renderLogo(logo, contentWidth))
	} else {
		b.WriteString(m.styles.header.Render("WeazlTunes"))
	}
	b.WriteString("\n\n")
	help := "[1] presets  [2] SomaFM  [3] Xiph  [4] my stations  [/] tune  [brackets] move  [ctrl+p] preset  [ctrl+r] rename  [ctrl+d] delete  [space] pause  [enter] play/add  [s] stop  [q] quit"
	b.WriteString(m.styles.header.Render(ansi.Wordwrap(help, contentWidth, " []/")))
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
		state := "now: "
		if m.paused {
			state = "paused: "
		}
		meter := " meter:fallback"
		if m.energy.Live {
			meter = " meter:live"
		}
		return m.styles.status.Render(ansi.Wordwrap(state+m.playing.Name+meter+" <"+m.playing.URL+">", max(20, m.width-4), " /_-<>"))
	}
	return m.styles.status.Render(ansi.Wordwrap(m.status, max(20, m.width-4), " /_-"))
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
