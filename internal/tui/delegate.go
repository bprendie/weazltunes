package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type stationDelegate struct {
	styles styles
}

func newStationDelegate() stationDelegate {
	return stationDelegate{styles: newStyles()}
}

func (d stationDelegate) Height() int  { return 4 }
func (d stationDelegate) Spacing() int { return 1 }
func (d stationDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d stationDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	st, ok := item.(stationItem)
	if !ok || m.Width() <= 0 {
		return
	}
	width := max(10, m.Width()-4)
	title := ansi.Truncate(st.Title(), width, "...")
	desc := ansi.Wordwrap(st.Description(), width, " /_-")
	desc = firstLines(desc, 2)
	titleStyle := d.styles.item
	descStyle := d.styles.help
	if index == m.Index() {
		titleStyle = d.styles.selected
		descStyle = lipgloss.NewStyle().Foreground(crushPurple)
		title = "> " + title
	} else {
		title = "  " + title
	}
	fmt.Fprintf(w, "%s\n%s", titleStyle.Render(title), descStyle.Render("  "+desc))
}

func firstLines(s string, limit int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= limit {
		return s
	}
	return strings.Join(lines[:limit], "\n")
}
