package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/sky0621/koro/internal/koro"
)

// DrawPlayer renders Koro as a hovering capsule drone with thrusters.
func DrawPlayer(dst *ebiten.Image, x, y, size float64, dir koro.Direction, body, accent color.Color) {
	cx := x + size/2
	cy := y + size/2
	capHeight := size * 0.65
	capRadius := capHeight / 2
	capWidth := size * 0.95
	left := cx - capWidth/2
	right := cx + capWidth/2
	top := cy - capHeight/2

	// Capsule body (rounded rectangle).
	FillRect(dst, left+capRadius, top, capWidth-capHeight, capHeight, body)
	FillCircle(dst, left+capRadius, cy, capRadius, body)
	FillCircle(dst, right-capRadius, cy, capRadius, body)

	// Accent bands.
	FillRect(dst, left+capRadius/2, top+capHeight*0.2, capWidth-capRadius, capHeight*0.15, accent)
	FillRect(dst, left+capRadius/3, top+capHeight*0.65, capWidth-capRadius/1.5, capHeight*0.12, accent)

	// Visor window.
	visorWidth := capWidth * 0.55
	visorHeight := capHeight * 0.45
	visorX := cx - visorWidth/2
	visorY := cy - visorHeight/2
	visorColor := color.RGBA{170, 245, 255, 235}
	FillRect(dst, visorX, visorY, visorWidth, visorHeight, visorColor)
	FillCircle(dst, visorX, visorY+visorHeight/2, visorHeight/2, visorColor)
	FillCircle(dst, visorX+visorWidth, visorY+visorHeight/2, visorHeight/2, visorColor)

	// Pupils indicate facing direction.
	eyeRadius := visorHeight * 0.25
	dirAngle := directionAngle(dir)
	offsetX := math.Cos(dirAngle) * visorWidth * 0.15
	offsetY := math.Sin(dirAngle) * visorHeight * 0.2
	pupilColor := color.RGBA{30, 60, 130, 255}
	FillCircle(dst, cx-visorWidth*0.15+offsetX, cy+offsetY, eyeRadius, pupilColor)
	FillCircle(dst, cx+visorWidth*0.15+offsetX, cy+offsetY, eyeRadius, pupilColor)

	drawThruster(dst, cx, cy, size, dir)
	drawFins(dst, cx, cy, size, accent)
}

// DrawGhost renders the enemies as angular shards with glowing cores.
func DrawGhost(dst *ebiten.Image, x, y, size float64, body color.Color, frightened bool) {
	cx := x + size/2
	cy := y + size/2

	top := y
	bottom := y + size
	left := x
	right := x + size

	coreColor := mixColor(body, color.RGBA{255, 255, 255, 80})
	if frightened {
		coreColor = mixColor(body, color.RGBA{255, 255, 255, 140})
	}

	FillTriangle(dst, cx, top, left+size*0.15, cy, right-size*0.15, cy, body)
	FillTriangle(dst, cx, bottom, left+size*0.15, cy, right-size*0.15, cy, body)
	FillTriangle(dst, cx, top+size*0.2, left+size*0.3, cy, right-size*0.3, cy, coreColor)
	FillTriangle(dst, cx, bottom-size*0.2, left+size*0.3, cy, right-size*0.3, cy, coreColor)

	// Vertical energy bands emphasise up/down motion.
	bandWidth := size * 0.12
	FillRect(dst, cx-bandWidth/2, y+size*0.1, bandWidth, size*0.8, mixColor(coreColor, color.RGBA{20, 20, 20, 120}))
	FillRect(dst, cx-bandWidth*1.6, y+size*0.25, bandWidth*0.6, size*0.5, mixColor(body, color.RGBA{0, 0, 0, 90}))
	FillRect(dst, cx+bandWidth, y+size*0.25, bandWidth*0.6, size*0.5, mixColor(body, color.RGBA{0, 0, 0, 90}))

	// Sensor eyes.
	eyeRadius := size * 0.12
	eyeOffset := size * 0.2
	eyeY := cy - size*0.05
	eyeColor := color.RGBA{255, 255, 255, 255}
	pupilColor := color.RGBA{20, 40, 90, 255}
	if frightened {
		pupilColor = color.RGBA{245, 245, 255, 210}
		eyeColor = color.RGBA{0, 0, 60, 255}
	}
	FillCircle(dst, cx-eyeOffset, eyeY, eyeRadius, eyeColor)
	FillCircle(dst, cx+eyeOffset, eyeY, eyeRadius, eyeColor)
	FillCircle(dst, cx-eyeOffset, eyeY+eyeRadius*0.3, eyeRadius*0.6, pupilColor)
	FillCircle(dst, cx+eyeOffset, eyeY+eyeRadius*0.3, eyeRadius*0.6, pupilColor)
}

func directionAngle(dir koro.Direction) float64 {
	switch dir {
	case koro.DirUp:
		return -math.Pi / 2
	case koro.DirDown:
		return math.Pi / 2
	case koro.DirLeft:
		return math.Pi
	case koro.DirRight:
		return 0
	default:
		return 0
	}
}

func drawThruster(dst *ebiten.Image, cx, cy, size float64, dir koro.Direction) {
	if dir == koro.DirNone {
		return
	}
	flameColor := color.RGBA{255, 140, 80, 230}
	switch dir {
	case koro.DirRight:
		FillTriangle(dst, cx-size*0.45, cy, cx-size*0.7, cy-size*0.15, cx-size*0.7, cy+size*0.15, flameColor)
	case koro.DirLeft:
		FillTriangle(dst, cx+size*0.45, cy, cx+size*0.7, cy-size*0.15, cx+size*0.7, cy+size*0.15, flameColor)
	case koro.DirUp:
		FillTriangle(dst, cx, cy+size*0.45, cx-size*0.15, cy+size*0.75, cx+size*0.15, cy+size*0.75, flameColor)
	case koro.DirDown:
		FillTriangle(dst, cx, cy-size*0.45, cx-size*0.15, cy-size*0.75, cx+size*0.15, cy-size*0.75, flameColor)
	}
}

func drawFins(dst *ebiten.Image, cx, cy, size float64, accent color.Color) {
	finColor := mixColor(accent, color.RGBA{255, 255, 255, 120})
	finWidth := size * 0.15
	finHeight := size * 0.25
	FillRect(dst, cx-size*0.35, cy+size*0.15, finWidth, finHeight, finColor)
	FillRect(dst, cx+size*0.2, cy+size*0.15, finWidth, finHeight, finColor)
	FillRect(dst, cx-size*0.35, cy-size*0.4, finWidth*0.8, finHeight*0.6, finColor)
	FillRect(dst, cx+size*0.25, cy-size*0.4, finWidth*0.8, finHeight*0.6, finColor)
}

func mixColor(a, b color.Color) color.Color {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return color.RGBA{
		R: uint8(((ar*3/4 + br/4) >> 8) & 0xff),
		G: uint8(((ag*3/4 + bg/4) >> 8) & 0xff),
		B: uint8(((ab*3/4 + bb/4) >> 8) & 0xff),
		A: uint8(((aa + ba/2) >> 8) & 0xff),
	}
}
