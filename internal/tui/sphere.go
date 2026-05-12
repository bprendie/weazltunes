package tui

import (
	"math"
	"strings"

	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

type SphereVisualizer struct {
	spring harmonica.Spring
	smile  motion
	blink  motion
	bob    motion
	tick   int
}

type motion struct {
	value    float64
	velocity float64
}

func NewSphereVisualizer(delta float64) SphereVisualizer {
	return SphereVisualizer{
		spring: harmonica.NewSpring(delta, 10.5, 0.42),
		smile:  motion{value: 0.8},
	}
}

func (s *SphereVisualizer) Step(playing bool) {
	s.tick++
	energy := 0.35
	if playing {
		energy = 1
	}
	s.smile.update(s.spring, 0.55+0.35*math.Sin(float64(s.tick)*0.08)*energy)
	s.bob.update(s.spring, math.Sin(float64(s.tick)*0.11)*energy)
	blinkTarget := 0.0
	if s.tick%130 > 118 {
		blinkTarget = 1
	}
	s.blink.update(s.spring, blinkTarget)
}

func (s SphereVisualizer) View(width int) string {
	w, h := sphereSize(width)
	var out strings.Builder
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			out.WriteString(s.cell(col, row, w, h))
		}
		if row < h-1 {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func (s SphereVisualizer) cell(col, row, w, h int) string {
	x := (float64(col)/float64(w-1))*2 - 1
	y := ((float64(row)/float64(h-1))*2 - 1) * 1.65
	y -= s.bob.value * 0.04
	d := x*x + y*y
	if d > 1.02 {
		return " "
	}
	if s.isEye(x, y, -0.34) || s.isEye(x, y, 0.34) || s.isMouth(x, y) {
		return lipgloss.NewStyle().Foreground(crushPurple).Bold(true).Render("█")
	}
	if d > 0.86 {
		return lipgloss.NewStyle().Foreground(crushPurple).Bold(true).Render("▓")
	}
	if s.isCheek(x, y) {
		return lipgloss.NewStyle().Foreground(crushPink).Bold(true).Render("▒")
	}
	if math.Sin((x+y+float64(s.tick)*0.03)*10) > 0.78 {
		return lipgloss.NewStyle().Foreground(crushMint).Render("░")
	}
	return lipgloss.NewStyle().Foreground(crushGold).Bold(true).Render("█")
}

func (s SphereVisualizer) isEye(x, y, cx float64) bool {
	eyeHeight := 0.11 * (1 - 0.78*s.blink.value)
	dx := (x - cx) / 0.12
	dy := (y + 0.24) / maxFloat(0.025, eyeHeight)
	return dx*dx+dy*dy < 1
}

func (s SphereVisualizer) isMouth(x, y float64) bool {
	if math.Abs(x) > 0.46 {
		return false
	}
	curve := 0.25 + s.smile.value*0.26*(x*x*2.8-0.35)
	return math.Abs(y-curve) < 0.055
}

func (s SphereVisualizer) isCheek(x, y float64) bool {
	return math.Abs(y-0.14) < 0.08 && (math.Abs(x-0.58) < 0.08 || math.Abs(x+0.58) < 0.08)
}

func (m *motion) update(spring harmonica.Spring, target float64) {
	m.value, m.velocity = spring.Update(m.value, m.velocity, target)
}
