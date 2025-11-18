package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/sky0621/koro/internal/koro"
)

type Game struct {
	k *koro.Koro
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.k.MoveLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.k.MoveRight()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.k.MoveUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.k.MoveDown()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Koroの仮の姿：白い四角
	ebitenutil.DrawRect(screen, g.k.X, g.k.Y, 20, 20, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 240, 320
}

func main() {
	g := &Game{
		k: koro.New(100, 100),
	}

	ebiten.SetWindowSize(480, 640)
	ebiten.SetWindowTitle("Koro Game")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
