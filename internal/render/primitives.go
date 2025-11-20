package render

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
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
