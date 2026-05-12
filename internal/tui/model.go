package tui

import (
	"fmt"
	"time"

	"github.com/bprendie/weazltunes/internal/audio"
	"github.com/bprendie/weazltunes/internal/config"
	"github.com/bprendie/weazltunes/internal/directory"
	"github.com/bprendie/weazltunes/internal/player"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
)

type mode int

const (
	modePresets mode = iota
	modeSoma
	modeXiph
	modeMyStations
)

const tunePrompt = "tune > "

type Model struct {
	cfg        config.Config
	styles     styles
	client     directory.Client
	player     *player.Player
	meter      *audio.Meter
	list       list.Model
	input      textinput.Model
	mode       mode
	status     string
	err        string
	playing    *directory.Station
	paused     bool
	energy     audio.Sample
	renaming   *directory.Station
	width      int
	height     int
	visualizer Visualizer
}

type stationItem struct {
	station directory.Station
}

func (i stationItem) Title() string { return i.station.Name }
func (i stationItem) Description() string {
	return fmt.Sprintf("%s  %s", i.station.Source, i.station.Description)
}
func (i stationItem) FilterValue() string {
	return i.station.Name + " " + i.station.Genre + " " + i.station.Description
}

type loadedMsg struct {
	stations []directory.Station
	status   string
}

type errMsg struct{ err error }
type tickMsg time.Time

func New(cfg config.Config) Model {
	input := newSearchInput()
	stations := newStationList()
	m := Model{
		cfg:        cfg,
		styles:     newStyles(),
		client:     directory.NewClient(),
		player:     player.New(),
		list:       stations,
		input:      input,
		mode:       modePresets,
		status:     "ready",
		visualizer: NewVisualizer(harmonica.FPS(30)),
	}
	m.setStations(presetStations(cfg))
	return m
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func newSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "genre, station search, or stream url"
	ti.Prompt = tunePrompt
	ti.CharLimit = 240
	ti.Width = 42
	return ti
}

func newStationList() list.Model {
	l := list.New(nil, newStationDelegate(), 80, 20)
	l.Title = "top 5 presets"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	return l
}
