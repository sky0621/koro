package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/sky0621/koro/internal/koro"
)

// Manager normalises keyboard, mouse, and touch inputs to a single direction.
type Manager struct {
	current        koro.Direction
	touchOrigins   map[ebiten.TouchID]touchPoint
	swipeThreshold float64
}

type touchPoint struct {
	x float64
	y float64
}

// NewManager creates an input manager with default thresholds.
func NewManager() *Manager {
	return &Manager{
		touchOrigins:   map[ebiten.TouchID]touchPoint{},
		swipeThreshold: 24,
	}
}

// Update samples the current input devices.
func (m *Manager) Update() {
	dir := m.keyboardDirection()
	if dir == koro.DirNone {
		dir = m.touchDirection()
	}
	if dir != koro.DirNone {
		m.current = dir
	}
}

// Direction returns the latest requested direction.
func (m *Manager) Direction() koro.Direction {
	return m.current
}

func (m *Manager) keyboardDirection() koro.Direction {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA):
		return koro.DirLeft
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD):
		return koro.DirRight
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW):
		return koro.DirUp
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS):
		return koro.DirDown
	default:
		return koro.DirNone
	}
}

func (m *Manager) touchDirection() koro.Direction {
	ids := ebiten.TouchIDs()
	if len(ids) == 0 {
		m.touchOrigins = map[ebiten.TouchID]touchPoint{}
		return koro.DirNone
	}

	seen := map[ebiten.TouchID]struct{}{}
	for _, id := range ids {
		seen[id] = struct{}{}
		x, y := ebiten.TouchPosition(id)
		pt := touchPoint{x: float64(x), y: float64(y)}
		if _, ok := m.touchOrigins[id]; !ok {
			m.touchOrigins[id] = pt
			continue
		}

		origin := m.touchOrigins[id]
		dx := pt.x - origin.x
		dy := pt.y - origin.y
		if math.Abs(dx) < m.swipeThreshold && math.Abs(dy) < m.swipeThreshold {
			continue
		}

		m.touchOrigins[id] = pt
		if math.Abs(dx) > math.Abs(dy) {
			if dx > 0 {
				m.cleanupTouches(seen)
				return koro.DirRight
			}
			m.cleanupTouches(seen)
			return koro.DirLeft
		}

		if dy > 0 {
			m.cleanupTouches(seen)
			return koro.DirDown
		}
		m.cleanupTouches(seen)
		return koro.DirUp
	}

	m.cleanupTouches(seen)
	return koro.DirNone
}

func (m *Manager) cleanupTouches(seen map[ebiten.TouchID]struct{}) {
	for id := range m.touchOrigins {
		if _, ok := seen[id]; !ok {
			delete(m.touchOrigins, id)
		}
	}
}
