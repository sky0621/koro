package ghost

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/sky0621/koro/internal/koro"
	"github.com/sky0621/koro/internal/level"
)

const (
	randomChangeChance   = 0.4
	randomChangeFrighten = 0.2
	visitPenaltyWeight   = 5.0
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
	visited         map[level.GridPos]int
}

// New creates a new ghost positioned at (x, y).
func New(x, y, tileSize float64, clr color.Color) *Ghost {
	body := koro.New(x, y, tileSize)
	body.SetSpeed(1.35)
	g := &Ghost{
		body:         body,
		baseSpeed:    body.Speed,
		primaryColor: clr,
		spawnX:       x,
		spawnY:       y,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
		visited:      map[level.GridPos]int{},
	}
	g.RespawnAt(x, y)
	return g
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
	} else if current := g.body.Direction(); current != koro.DirNone {
		g.body.SetIntentDirection(current)
	}

	g.body.Update(l)
	g.recordVisit(l)
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
	g.RespawnAt(g.spawnX, g.spawnY)
}

// RespawnAt teleports the ghost to a new spawn position.
func (g *Ghost) RespawnAt(x, y float64) {
	g.spawnX = x
	g.spawnY = y
	g.body.SetPosition(x, y)
	g.body.SetIntentDirection(koro.DirNone)
	g.frightenedTimer = 0
	g.visited = map[level.GridPos]int{}
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
	current := g.body.Direction()
	changeChance := randomChangeChance
	if g.IsFrightened() {
		changeChance = randomChangeFrighten
	}
	randomRecalc := changeChance > 0 && g.rng.Float64() < changeChance
	needDecision := current == koro.DirNone || !g.body.CanMove(l, current) || g.atIntersection(l) || randomRecalc
	if !needDecision {
		return koro.DirNone
	}

	options := g.availableDirections(l)
	if len(options) == 0 {
		if current != koro.DirNone && g.body.CanMove(l, current) {
			return current
		}
		return koro.DirNone
	}

	g.shuffleDirections(options)

	if g.IsFrightened() {
		return options[g.rng.Intn(len(options))]
	}

	cx, cy := g.body.Center()
	tileSize := float64(l.TileSize)
	bestDir := options[0]
	bestScore := math.MaxFloat64

	for _, dir := range options {
		dx, dy := dir.Delta()
		nextX := cx + float64(dx)*tileSize
		nextY := cy + float64(dy)*tileSize
		dist := math.Hypot(targetX-nextX, targetY-nextY)
		grid := g.gridAhead(l, dir)
		visitScore := float64(g.visited[grid]) * visitPenaltyWeight
		score := dist + visitScore
		if score < bestScore {
			bestScore = score
			bestDir = dir
		}
	}

	return bestDir
}

func (g *Ghost) shuffleDirections(dirs []koro.Direction) {
	g.rng.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})
}

func (g *Ghost) atIntersection(l *level.Level) bool {
	cx, cy := g.body.Center()
	grid := l.GridForPixel(cx, cy)
	tile := float64(l.TileSize)
	centerX := float64(grid.Col)*tile + tile/2
	centerY := float64(grid.Row)*tile + tile/2
	tolerance := tile * 0.3
	return math.Abs(centerX-cx) < tolerance && math.Abs(centerY-cy) < tolerance
}

func (g *Ghost) availableDirections(l *level.Level) []koro.Direction {
	current := g.body.Direction()
	opposite := oppositeDirection(current)

	dirs := []koro.Direction{koro.DirUp, koro.DirDown, koro.DirLeft, koro.DirRight}
	valid := dirs[:0]
	for _, dir := range dirs {
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

func (g *Ghost) recordVisit(l *level.Level) {
	cx, cy := g.body.Center()
	grid := l.GridForPixel(cx, cy)
	g.visited[grid]++
}

func (g *Ghost) gridAhead(l *level.Level, dir koro.Direction) level.GridPos {
	cx, cy := g.body.Center()
	dx, dy := dir.Delta()
	nextX := cx + float64(dx)*float64(l.TileSize)
	nextY := cy + float64(dy)*float64(l.TileSize)
	return l.GridForPixel(nextX, nextY)
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
