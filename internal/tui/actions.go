package tui

import (
	"github.com/bprendie/weazltunes/internal/directory"
	"github.com/charmbracelet/bubbles/list"
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

func (m *Model) setStations(stations []directory.Station) {
	items := make([]list.Item, 0, len(stations))
	for _, st := range stations {
		items = append(items, stationItem{station: st})
	}
	m.list.SetItems(items)
}

func directoryStation(item stationItem) directory.Station {
	return item.station
}
