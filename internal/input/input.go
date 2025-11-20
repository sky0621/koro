package input

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/sky0621/koro/internal/koro"
)

// Manager normalises keyboard/gamepad input to a single direction.
type Manager struct {
	current koro.Direction
}

// NewManager creates an input manager with default thresholds.
func NewManager() *Manager {
	return &Manager{}
}

// Update samples the current input devices.
func (m *Manager) Update() {
	dir := m.keyboardDirection()
	if dir == koro.DirNone {
		dir = m.gamepadDirection()
	}
	m.current = dir
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

func (m *Manager) gamepadDirection() koro.Direction {
	for _, id := range ebiten.GamepadIDs() {
		if !ebiten.IsStandardGamepadLayoutAvailable(id) {
			continue
		}
		switch {
		case ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftLeft):
			return koro.DirLeft
		case ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftRight):
			return koro.DirRight
		case ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftTop):
			return koro.DirUp
		case ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftBottom):
			return koro.DirDown
		}
	}
	return koro.DirNone
}
