package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/sky0621/koro/internal/koro"
	"github.com/sky0621/koro/internal/level"
)

type Game struct {
	k     *koro.Koro
	level *level.Level
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.k.QueueDirection(koro.DirLeft)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.k.QueueDirection(koro.DirRight)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.k.QueueDirection(koro.DirUp)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.k.QueueDirection(koro.DirDown)
	}

	g.k.Update(g.level)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawLevel(screen)
	drawSolidRect(screen, g.k.X, g.k.Y, g.k.Size, g.k.Size, color.RGBA{255, 255, 0, 255})
}

func (g *Game) drawLevel(screen *ebiten.Image) {
	tileSize := float64(g.level.TileSize)

	for row := 0; row < g.level.Height; row++ {
		for col := 0; col < g.level.Width; col++ {
			x := float64(col) * tileSize
			y := float64(row) * tileSize

			var c color.Color
			switch g.level.TileAt(col, row) {
			case level.TileWall:
				c = color.NRGBA{0, 0, 80, 255}
			case level.TileWarp:
				c = color.NRGBA{30, 30, 90, 255}
			default:
				c = color.NRGBA{10, 10, 10, 255}
			}

			drawSolidRect(screen, x, y, tileSize, tileSize, c)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.level.PixelWidth(), g.level.PixelHeight()
}

func main() {
	lvl := level.DefaultLevel()
	tileSize := float64(lvl.TileSize)

	g := &Game{
		level: lvl,
		k:     koro.New(7*tileSize, 11*tileSize, tileSize),
	}

	ebiten.SetWindowSize(lvl.PixelWidth()*2, lvl.PixelHeight()*2)
	ebiten.SetWindowTitle("Koro Game")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

var whitePixel = ebiten.NewImage(1, 1)

func init() {
	whitePixel.Fill(color.White)
}

func drawSolidRect(dst *ebiten.Image, x, y, width, height float64, clr color.Color) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(width, height)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	dst.DrawImage(whitePixel, op)
}
