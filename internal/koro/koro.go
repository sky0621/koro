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

func (d Direction) Delta() (int, int) {
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
	X, Y   float64
	Size   float64
	Speed  float64
	dir    Direction
	intent Direction
}

// New returns a configured Koro instance positioned at x,y with the provided size.
func New(x, y, size float64) *Koro {
	return &Koro{
		X:    x,
		Y:    y,
		Size: size,
		// Speed is tuned for smooth per-frame pixel movement.
		Speed:  1.5,
		dir:    DirNone,
		intent: DirNone,
	}
}

// SetIntentDirection stores the desired direction from player input.
func (k *Koro) SetIntentDirection(dir Direction) {
	k.intent = dir
}

// Update moves Koro according to queued directions and level collisions.
func (k *Koro) Update(l *level.Level) {
	tileSize := float64(l.TileSize)
	k.applyIntent(l, tileSize)

	if k.dir == DirNone {
		return
	}

	if k.canMove(l, k.dir) {
		dx, dy := k.dir.Delta()
		k.X += float64(dx) * k.Speed
		k.Y += float64(dy) * k.Speed
		k.handleWarp(l)
		return
	}

	k.snapToGrid(tileSize)
	k.dir = DirNone
}

func (k *Koro) applyIntent(l *level.Level, tileSize float64) {
	if k.intent == DirNone {
		k.dir = DirNone
		return
	}

	if k.intent == k.dir {
		return
	}

	if !k.canMove(l, k.intent) {
		return
	}

	k.snapAxisForDirection(tileSize, k.intent)
	k.dir = k.intent
}

func (k *Koro) canMove(l *level.Level, dir Direction) bool {
	dx, dy := dir.Delta()
	nextX := k.X + float64(dx)*k.Speed
	nextY := k.Y + float64(dy)*k.Speed
	return !l.Collides(nextX, nextY, k.Size)
}

// CanMove reports whether moving in the provided direction would collide.
func (k *Koro) CanMove(l *level.Level, dir Direction) bool {
	return k.canMove(l, dir)
}

func (k *Koro) snapAxisForDirection(tileSize float64, dir Direction) {
	switch dir {
	case DirUp, DirDown:
		k.X = math.Round(k.X/tileSize) * tileSize
	case DirLeft, DirRight:
		k.Y = math.Round(k.Y/tileSize) * tileSize
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

// Direction returns the current heading.
func (k *Koro) Direction() Direction {
	return k.dir
}

// SetPosition teleports Koro to the provided coordinates and clears direction state.
func (k *Koro) SetPosition(x, y float64) {
	k.X = x
	k.Y = y
	k.dir = DirNone
	k.intent = DirNone
}

// Center returns the current center point coordinates.
func (k *Koro) Center() (float64, float64) {
	return k.X + k.Size/2, k.Y + k.Size/2
}

// SetSpeed adjusts movement speed.
func (k *Koro) SetSpeed(speed float64) {
	k.Speed = speed
}
