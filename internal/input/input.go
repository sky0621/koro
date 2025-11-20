package input

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/sky0621/koro/internal/koro"
)

// Manager normalises keyboard input to a single direction.
type Manager struct {
	current koro.Direction
}

// NewManager creates an input manager with default thresholds.
func NewManager() *Manager {
	return &Manager{}
}

// Update samples the current keyboard input.
func (m *Manager) Update() {
	if dir := m.keyboardDirection(); dir != koro.DirNone {
		m.current = dir
	}
}

// Direction returns the latest requested direction.
func (m *Manager) Direction() koro.Direction {
	return m.current
}

func (m *Manager) keyboardDirection() koro.Direction {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyArrowLeft):
		return koro.DirLeft
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight):
		return koro.DirRight
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		return koro.DirUp
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		return koro.DirDown
	default:
		return koro.DirNone
	}
}
