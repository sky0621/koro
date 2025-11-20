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
		dx, dy := k.dir.Delta()
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
	dx, dy := dir.Delta()
	nextX := k.X + float64(dx)*k.Speed
	nextY := k.Y + float64(dy)*k.Speed
	return !l.Collides(nextX, nextY, k.Size)
}

// CanMove reports whether moving in the provided direction would collide.
func (k *Koro) CanMove(l *level.Level, dir Direction) bool {
	return k.canMove(l, dir)
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

// Direction returns the current heading.
func (k *Koro) Direction() Direction {
	return k.dir
}

// SetPosition teleports Koro to the provided coordinates and clears direction state.
func (k *Koro) SetPosition(x, y float64) {
	k.X = x
	k.Y = y
	k.dir = DirNone
	k.next = DirNone
}

// Center returns the current center point coordinates.
func (k *Koro) Center() (float64, float64) {
	return k.X + k.Size/2, k.Y + k.Size/2
}

// SetSpeed adjusts movement speed.
func (k *Koro) SetSpeed(speed float64) {
	k.Speed = speed
}
