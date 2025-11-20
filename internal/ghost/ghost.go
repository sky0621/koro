package ghost

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/sky0621/koro/internal/koro"
	"github.com/sky0621/koro/internal/level"
)

// Ghost encapsulates enemy behaviour with simple chase logic.
type Ghost struct {
	body            *koro.Koro
	baseSpeed       float64
	primaryColor    color.Color
	frightenedTimer int
	rng             *rand.Rand
	spawnX          float64
	spawnY          float64
}

// New creates a new ghost positioned at (x, y).
func New(x, y, tileSize float64, clr color.Color) *Ghost {
	body := koro.New(x, y, tileSize)
	body.SetSpeed(1.35)
	return &Ghost{
		body:         body,
		baseSpeed:    body.Speed,
		primaryColor: clr,
		spawnX:       x,
		spawnY:       y,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Update advances the ghost AI and movement.
func (g *Ghost) Update(l *level.Level, targetX, targetY float64) {
	if g.frightenedTimer > 0 {
		g.frightenedTimer--
	}

	if g.IsFrightened() {
		g.body.SetSpeed(g.baseSpeed * 0.75)
	} else {
		g.body.SetSpeed(g.baseSpeed)
	}

	if dir := g.nextDirection(l, targetX, targetY); dir != koro.DirNone {
		g.body.SetIntentDirection(dir)
	}

	g.body.Update(l)
}

// Color returns the draw color according to current state.
func (g *Ghost) Color() color.Color {
	if g.IsFrightened() {
		return color.RGBA{0, 0, 255, 255}
	}
	return g.primaryColor
}

// SetFrightened activates frightened mode for the provided frame duration.
func (g *Ghost) SetFrightened(duration int) {
	if duration > g.frightenedTimer {
		g.frightenedTimer = duration
	}
}

// IsFrightened reports whether the ghost is currently vulnerable.
func (g *Ghost) IsFrightened() bool {
	return g.frightenedTimer > 0
}

// Reset moves the ghost back to its spawn point.
func (g *Ghost) Reset() {
	g.body.SetPosition(g.spawnX, g.spawnY)
	g.frightenedTimer = 0
}

// Position returns the current location.
func (g *Ghost) Position() (float64, float64) {
	return g.body.X, g.body.Y
}

// Size returns the body size in pixels.
func (g *Ghost) Size() float64 {
	return g.body.Size
}

// Body returns the internal mover component.
func (g *Ghost) Body() *koro.Koro {
	return g.body
}

func (g *Ghost) nextDirection(l *level.Level, targetX, targetY float64) koro.Direction {
	if !g.atIntersection(l) {
		return koro.DirNone
	}

	options := g.availableDirections(l)
	if len(options) == 0 {
		return koro.DirNone
	}

	if g.IsFrightened() {
		return options[g.rng.Intn(len(options))]
	}

	cx, cy := g.body.Center()
	tileSize := float64(l.TileSize)
	bestDir := options[0]
	bestDist := math.MaxFloat64

	for _, dir := range options {
		dx, dy := dir.Delta()
		nextX := cx + float64(dx)*tileSize
		nextY := cy + float64(dy)*tileSize
		dist := math.Hypot(targetX-nextX, targetY-nextY)
		if dist < bestDist {
			bestDist = dist
			bestDir = dir
		}
	}

	if len(options) > 1 && g.rng.Float64() < 0.2 {
		return options[g.rng.Intn(len(options))]
	}

	return bestDir
}

func (g *Ghost) atIntersection(l *level.Level) bool {
	cx, cy := g.body.Center()
	grid := l.GridForPixel(cx, cy)
	tile := float64(l.TileSize)
	centerX := float64(grid.Col)*tile + tile/2
	centerY := float64(grid.Row)*tile + tile/2
	const tolerance = 0.2
	return math.Abs(centerX-cx) < tolerance && math.Abs(centerY-cy) < tolerance
}

func (g *Ghost) availableDirections(l *level.Level) []koro.Direction {
	all := []koro.Direction{koro.DirUp, koro.DirDown, koro.DirLeft, koro.DirRight}
	current := g.body.Direction()
	opposite := oppositeDirection(current)
	valid := make([]koro.Direction, 0, len(all))
	for _, dir := range all {
		if dir == opposite {
			continue
		}
		if g.body.CanMove(l, dir) {
			valid = append(valid, dir)
		}
	}
	if len(valid) == 0 && opposite != koro.DirNone && g.body.CanMove(l, opposite) {
		return []koro.Direction{opposite}
	}
	return valid
}

func oppositeDirection(dir koro.Direction) koro.Direction {
	switch dir {
	case koro.DirLeft:
		return koro.DirRight
	case koro.DirRight:
		return koro.DirLeft
	case koro.DirUp:
		return koro.DirDown
	case koro.DirDown:
		return koro.DirUp
	default:
		return koro.DirNone
	}
}
