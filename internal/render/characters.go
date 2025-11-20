package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/sky0621/koro/internal/koro"
)

// DrawPlayer renders a wedge-mouthed hero facing the provided direction.
func DrawPlayer(dst *ebiten.Image, x, y, size float64, dir koro.Direction, body, mouth color.Color) {
	radius := size / 2
	cx := x + radius
	cy := y + radius
	FillCircle(dst, cx, cy, radius, body)

	angle := directionAngle(dir)
	mouthOpen := math.Pi / 6
	ax := cx
	ay := cy
	bx := cx + math.Cos(angle+mouthOpen)*radius
	by := cy + math.Sin(angle+mouthOpen)*radius
	cx2 := cx + math.Cos(angle-mouthOpen)*radius
	cy2 := cy + math.Sin(angle-mouthOpen)*radius
	FillTriangle(dst, ax, ay, bx, by, cx2, cy2, mouth)
}

// DrawGhost renders a rounded ghost shape with eyes.
func DrawGhost(dst *ebiten.Image, x, y, size float64, body color.Color, frightened bool) {
	radius := size / 2
	cx := x + radius
	cy := y + radius

	FillCircle(dst, cx, cy, radius, body)
	FillRect(dst, x, y+radius*0.6, size, radius*0.8, body)

	footRadius := size / 6
	for i := 0; i < 3; i++ {
		footCenterX := x + footRadius + float64(i)*(2*footRadius)
		footCenterY := y + size - footRadius
		FillCircle(dst, footCenterX, footCenterY, footRadius, body)
	}

	eyeY := y + radius*0.7
	eyeOffset := radius * 0.4
	eyeRadius := radius * 0.22
	pupilRadius := eyeRadius * 0.4

	var eyeColor color.Color = color.White
	var pupilColor color.Color = color.RGBA{10, 30, 120, 255}
	if frightened {
		pupilColor = color.RGBA{255, 255, 255, 255}
		eyeColor = color.RGBA{0, 0, 50, 255}
	}

	leftEyeX := cx - eyeOffset
	rightEyeX := cx + eyeOffset

	FillCircle(dst, leftEyeX, eyeY, eyeRadius, eyeColor)
	FillCircle(dst, rightEyeX, eyeY, eyeRadius, eyeColor)

	pupilOffsetY := radius * 0.15
	if frightened {
		pupilOffsetY = 0
	}

	FillCircle(dst, leftEyeX, eyeY+pupilOffsetY, pupilRadius, pupilColor)
	FillCircle(dst, rightEyeX, eyeY+pupilOffsetY, pupilRadius, pupilColor)
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
