package render

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whitePixelOnce sync.Once
	whitePixel     *ebiten.Image
)

func ensureWhitePixel() {
	whitePixelOnce.Do(func() {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	})
}

// FillRect draws a solid rectangle tinted with the provided color.
func FillRect(dst *ebiten.Image, x, y, width, height float64, clr color.Color) {
	if width <= 0 || height <= 0 {
		return
	}
	ensureWhitePixel()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(width, height)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	dst.DrawImage(whitePixel, op)
}

// FillCircle draws a filled circle at the provided center and radius.
func FillCircle(dst *ebiten.Image, cx, cy, radius float64, clr color.Color) {
	if radius <= 0 {
		return
	}
	vector.FillCircle(dst, float32(cx), float32(cy), float32(radius), clr, false)
}

// FillTriangle draws a triangle using the three provided points.
func FillTriangle(dst *ebiten.Image, ax, ay, bx, by, cx, cy float64, clr color.Color) {
	var path vector.Path
	path.MoveTo(float32(ax), float32(ay))
	path.LineTo(float32(bx), float32(by))
	path.LineTo(float32(cx), float32(cy))
	path.Close()
	drawOpts := &vector.DrawPathOptions{}
	drawOpts.ColorScale.ScaleWithColor(clr)
	vector.FillPath(dst, &path, nil, drawOpts)
}
