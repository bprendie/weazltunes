package tui

func sphereSize(width int) (int, int) {
	w := width - 4
	if w > 52 {
		w = 52
	}
	if w < 22 {
		w = 22
	}
	h := w / 3
	if h > 17 {
		h = 17
	}
	if h < 8 {
		h = 8
	}
	return w, h
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
