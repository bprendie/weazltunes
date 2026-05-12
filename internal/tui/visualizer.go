package tui

import (
	"math"
	"strings"

	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

type visualizerMode int

const (
	visualizerSphere visualizerMode = iota
	visualizerBars
)

type Visualizer struct {
	mode       visualizerMode
	spring     harmonica.Spring
	bars       []float64
	velocities []float64
	sphere     SphereVisualizer
	tick       int
}

func NewVisualizer(delta float64) Visualizer {
	return Visualizer{
		mode:       visualizerSphere,
		spring:     harmonica.NewSpring(delta, 9.0, 0.35),
		bars:       make([]float64, 24),
		velocities: make([]float64, 24),
		sphere:     NewSphereVisualizer(delta),
	}
}

func (v *Visualizer) Step(playing bool) {
	v.tick++
	v.stepBars(playing)
	v.sphere.Step(playing)
}

func (v *Visualizer) Toggle() string {
	if v.mode == visualizerSphere {
		v.mode = visualizerBars
		return "bars"
	}
	v.mode = visualizerSphere
	return "sphere"
}

func (v Visualizer) View(styles styles, width int) string {
	if v.mode == visualizerBars {
		return v.barsView()
	}
	return v.sphere.View(width)
}

func (v *Visualizer) stepBars(playing bool) {
	for i := range v.bars {
		base := 2.0
		if playing {
			base = 4 + 10*(0.5+0.5*math.Sin(float64(v.tick+i)*0.35))
		}
		target := base + 5*(0.5+0.5*math.Sin(float64(v.tick)*0.12+float64(i)*0.9))
		v.bars[i], v.velocities[i] = v.spring.Update(v.bars[i], v.velocities[i], target)
	}
}

func (v Visualizer) barsView() string {
	var b strings.Builder
	blocks := []rune("▁▂▃▄▅▆▇█")
	for i, value := range v.bars {
		b.WriteString(lipgloss.NewStyle().Foreground(v.color(i)).Render(string(blocks[v.index(value)])))
		b.WriteRune(' ')
	}
	return b.String()
}

func (v Visualizer) index(value float64) int {
	idx := int(math.Round(value / 2))
	if idx < 0 {
		return 0
	}
	if idx > 7 {
		return 7
	}
	return idx
}

func (v Visualizer) color(i int) lipgloss.Color {
	if i%3 == 0 {
		return crushPink
	}
	if i%3 == 1 {
		return crushMint
	}
	return crushPurple
}
