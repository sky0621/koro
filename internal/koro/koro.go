package koro

import (
	"math"

	"github.com/sky0621/koro/internal/level"
)

// Direction expresses which way Koro is heading.
type Direction int

const (
	DirNone Direction = iota
	DirLeft
	DirRight
	DirUp
	DirDown
)

func (d Direction) delta() (int, int) {
	switch d {
	case DirLeft:
		return -1, 0
	case DirRight:
		return 1, 0
	case DirUp:
		return 0, -1
	case DirDown:
		return 0, 1
	default:
		return 0, 0
	}
}

// Koro represents the controllable hero.
type Koro struct {
	X, Y  float64
	Size  float64
	Speed float64
	dir   Direction
	next  Direction
}

// New returns a configured Koro instance positioned at x,y with the provided size.
func New(x, y, size float64) *Koro {
	return &Koro{
		X:    x,
		Y:    y,
		Size: size,
		// Speed is tuned for smooth per-frame pixel movement.
		Speed: 1.5,
		dir:   DirNone,
		next:  DirNone,
	}
}

// QueueDirection stores the next intended direction from player input.
func (k *Koro) QueueDirection(dir Direction) {
	k.next = dir
}

// Update moves Koro according to queued directions and level collisions.
func (k *Koro) Update(l *level.Level) {
	tileSize := float64(l.TileSize)
	k.tryApplyQueuedDirection(l, tileSize)

	if k.dir == DirNone {
		return
	}

	if k.canMove(l, k.dir) {
		dx, dy := k.dir.delta()
		k.X += float64(dx) * k.Speed
		k.Y += float64(dy) * k.Speed
		k.handleWarp(l)
		return
	}

	k.snapToGrid(tileSize)
	k.dir = DirNone
}

func (k *Koro) tryApplyQueuedDirection(l *level.Level, tileSize float64) {
	if k.next == DirNone {
		return
	}

	if !k.isAlignedFor(k.next, tileSize) {
		return
	}

	if !k.canMove(l, k.next) {
		return
	}

	k.snapToGrid(tileSize)
	k.dir = k.next
	k.next = DirNone
}

func (k *Koro) canMove(l *level.Level, dir Direction) bool {
	dx, dy := dir.delta()
	nextX := k.X + float64(dx)*k.Speed
	nextY := k.Y + float64(dy)*k.Speed
	return !l.Collides(nextX, nextY, k.Size)
}

func (k *Koro) isAlignedFor(dir Direction, tileSize float64) bool {
	const minThreshold = 0.8
	threshold := math.Max(minThreshold, tileSize*0.2)
	if threshold > tileSize/2 {
		threshold = tileSize / 2
	}

	remX := math.Mod(k.X, tileSize)
	remY := math.Mod(k.Y, tileSize)
	if remX < 0 {
		remX += tileSize
	}
	if remY < 0 {
		remY += tileSize
	}

	switch dir {
	case DirUp, DirDown:
		return remX < threshold || tileSize-remX < threshold
	case DirLeft, DirRight:
		return remY < threshold || tileSize-remY < threshold
	default:
		return true
	}
}

func (k *Koro) snapToGrid(tileSize float64) {
	k.X = math.Round(k.X/tileSize) * tileSize
	k.Y = math.Round(k.Y/tileSize) * tileSize
}

func (k *Koro) handleWarp(l *level.Level) {
	centerX := k.X + k.Size/2
	centerY := k.Y + k.Size/2
	grid := l.GridForPixel(centerX, centerY)
	target, ok := l.WarpTarget(grid)
	if !ok {
		return
	}

	k.X = float64(target.Col * l.TileSize)
	k.Y = float64(target.Row * l.TileSize)
}
