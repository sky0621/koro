package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/sky0621/koro/internal/ghost"
	"github.com/sky0621/koro/internal/input"
	"github.com/sky0621/koro/internal/koro"
	"github.com/sky0621/koro/internal/level"
	"github.com/sky0621/koro/internal/render"
)

type Game struct {
	level    *level.Level
	tileSize float64

	player *koro.Koro
	ghosts []*ghost.Ghost
	input  *input.Manager

	score      int
	lives      int
	state      GameState
	readyTimer int
	powerTimer int

	playerSpawnX float64
	playerSpawnY float64
}

type GameState int

const (
	StateReady GameState = iota
	StatePlaying
	StateCleared
	StateGameOver
)

const (
	startLives         = 3
	pelletScore        = 10
	powerPelletScore   = 50
	ghostScore         = 200
	powerModeDuration  = 600
	readyDelayFrames   = 60
	collisionShrinkage = 0.85
)

var (
	colorWall        = color.NRGBA{0, 0, 80, 255}
	colorWarp        = color.NRGBA{30, 30, 90, 255}
	colorFloor       = color.NRGBA{10, 10, 10, 255}
	colorPlayer      = color.RGBA{255, 255, 0, 255}
	colorPowerPellet = color.RGBA{255, 165, 0, 255}
)

var ghostSpawnTiles = []struct {
	col int
	row int
	clr color.Color
}{
	{7, 9, color.RGBA{255, 0, 0, 255}},
	{6, 9, color.RGBA{0, 255, 255, 255}},
	{8, 9, color.RGBA{255, 105, 180, 255}},
}

func newGame() *Game {
	lvl := level.DefaultLevel()
	g := &Game{
		level:      lvl,
		tileSize:   float64(lvl.TileSize),
		input:      input.NewManager(),
		lives:      startLives,
		score:      0,
		state:      StateReady,
		readyTimer: readyDelayFrames,
	}
	g.setupActors()
	return g
}

func (g *Game) setupActors() {
	g.tileSize = float64(g.level.TileSize)
	g.playerSpawnX = 7 * g.tileSize
	g.playerSpawnY = 11 * g.tileSize
	g.player = koro.New(g.playerSpawnX, g.playerSpawnY, g.tileSize)
	g.player.SetSpeed(1.6)
	g.ghosts = make([]*ghost.Ghost, 0, len(ghostSpawnTiles))
	for _, spawn := range ghostSpawnTiles {
		x := float64(spawn.col) * g.tileSize
		y := float64(spawn.row) * g.tileSize
		gh := ghost.New(x, y, g.tileSize, spawn.clr)
		g.ghosts = append(g.ghosts, gh)
	}
	g.powerTimer = 0
}

func (g *Game) Update() error {
	switch g.state {
	case StateReady:
		if g.readyTimer > 0 {
			g.readyTimer--
			return nil
		}
		g.state = StatePlaying
	case StatePlaying:
		g.handleInput()
		g.player.Update(g.level)
		g.handlePelletPickup()
		g.updateGhosts()
		g.updatePowerTimer()
		g.checkGhostCollisions()
		if g.state == StateGameOver {
			return nil
		}
		if g.level.RemainingPellets() == 0 {
			g.state = StateCleared
			g.readyTimer = readyDelayFrames
		}
	case StateCleared:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.resetLevel(true)
		}
	case StateGameOver:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.resetLevel(false)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawLevel(screen)
	g.drawPellets(screen)
	g.drawGhosts(screen)
	render.DrawPlayer(screen, g.player.X, g.player.Y, g.player.Size, g.player.Direction(), colorPlayer, colorFloor)
	g.drawHUD(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.level.PixelWidth(), g.level.PixelHeight()
}

func (g *Game) drawLevel(screen *ebiten.Image) {
	tileSize := g.tileSize
	for row := 0; row < g.level.Height; row++ {
		for col := 0; col < g.level.Width; col++ {
			x := float64(col) * tileSize
			y := float64(row) * tileSize
			var c color.Color
			switch g.level.TileAt(col, row) {
			case level.TileWall:
				c = colorWall
			case level.TileWarp:
				c = colorWarp
			default:
				c = colorFloor
			}
			render.FillRect(screen, x, y, tileSize, tileSize, c)
		}
	}
}

func (g *Game) drawPellets(screen *ebiten.Image) {
	tileSize := g.tileSize
	half := tileSize / 2
	for row := 0; row < g.level.Height; row++ {
		for col := 0; col < g.level.Width; col++ {
			var size float64
			var clr color.Color
			switch g.level.PelletAt(col, row) {
			case level.PelletSmall:
				size = tileSize * 0.2
				clr = color.White
			case level.PelletPower:
				size = tileSize * 0.5
				clr = colorPowerPellet
			default:
				continue
			}
			centerX := float64(col)*tileSize + half
			centerY := float64(row)*tileSize + half
			render.FillRect(screen, centerX-size/2, centerY-size/2, size, size, clr)
		}
	}
}

func (g *Game) drawGhosts(screen *ebiten.Image) {
	for _, gh := range g.ghosts {
		x, y := gh.Position()
		render.DrawGhost(screen, x, y, gh.Size(), gh.Color(), gh.IsFrightened())
	}
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	text := fmt.Sprintf("Score: %d  Lives: %d", g.score, g.lives)
	switch g.state {
	case StateReady:
		text += "  Ready!"
	case StateCleared:
		text += "  LEVEL CLEAR - Press Enter"
	case StateGameOver:
		text += "  GAME OVER - Press Enter"
	}
	if g.powerTimer > 0 {
		text += fmt.Sprintf("  Power %ds", g.powerTimer/60)
	}
	ebitenutil.DebugPrint(screen, text)
}

func (g *Game) handleInput() {
	g.input.Update()
	if dir := g.input.Direction(); dir != koro.DirNone {
		g.player.QueueDirection(dir)
	}
}

func (g *Game) handlePelletPickup() {
	cx, cy := g.player.Center()
	grid := g.level.GridForPixel(cx, cy)
	switch g.level.ConsumePellet(grid.Col, grid.Row) {
	case level.PelletSmall:
		g.score += pelletScore
	case level.PelletPower:
		g.score += powerPelletScore
		g.activatePowerMode()
	}
}

func (g *Game) activatePowerMode() {
	g.powerTimer = powerModeDuration
	for _, gh := range g.ghosts {
		gh.SetFrightened(powerModeDuration)
	}
}

func (g *Game) updateGhosts() {
	px, py := g.player.Center()
	for _, gh := range g.ghosts {
		gh.Update(g.level, px, py)
	}
}

func (g *Game) updatePowerTimer() {
	if g.powerTimer > 0 {
		g.powerTimer--
	}
}

func (g *Game) checkGhostCollisions() {
	px, py := g.player.Center()
	playerRadius := g.player.Size / 2 * collisionShrinkage
	for _, gh := range g.ghosts {
		cx, cy := gh.Body().Center()
		ghostRadius := gh.Size() / 2 * collisionShrinkage
		if math.Hypot(px-cx, py-cy) <= playerRadius+ghostRadius {
			if gh.IsFrightened() {
				g.score += ghostScore
				gh.Reset()
			} else {
				g.loseLife()
			}
			return
		}
	}
}

func (g *Game) loseLife() {
	g.lives--
	if g.lives <= 0 {
		g.state = StateGameOver
		return
	}
	g.powerTimer = 0
	g.resetActorPositions()
	g.state = StateReady
	g.readyTimer = readyDelayFrames
}

func (g *Game) resetActorPositions() {
	g.player.SetPosition(g.playerSpawnX, g.playerSpawnY)
	for _, gh := range g.ghosts {
		gh.Reset()
	}
}

func (g *Game) resetLevel(keepScore bool) {
	g.level = level.DefaultLevel()
	g.setupActors()
	if !keepScore {
		g.score = 0
		g.lives = startLives
	}
	g.state = StateReady
	g.readyTimer = readyDelayFrames
}

func main() {
	g := newGame()
	ebiten.SetWindowSize(g.level.PixelWidth()*2, g.level.PixelHeight()*2)
	ebiten.SetWindowTitle("Koro Game")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
