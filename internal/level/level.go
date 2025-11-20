package level

import (
	"fmt"
	"math"
)

// TileType represents the type of a tile within the level map.
type TileType int

const (
	TileWall TileType = iota
	TilePath
	TileWarp
)

// GridPos holds column/row coordinates in the tile map.
type GridPos struct {
	Col int
	Row int
}

// Level contains tile data and warp links for a stage.
type Level struct {
	Tiles       [][]TileType
	TileSize    int
	Width       int
	Height      int
	warpTargets map[GridPos]GridPos
}

// DefaultLevel returns the built-in stage used for early development.
func DefaultLevel() *Level {
	layout := []string{
		"###############",
		"#.............#",
		"#.###.###.###.#",
		"#.#.........#.#",
		"#.#.###.###.#.#",
		"#.....#.......#",
		"###.#.#.#.#.###",
		"W...#.....#...W",
		"#.#.#.###.#.#.#",
		"#.#.#.....#.#.#",
		"#.#.#######.#.#",
		"#.............#",
		"###.###.#.###.#",
		"#.....#.#.....#",
		"#.###.#.#.###.#",
		"#.............#",
		"#.###########.#",
		"#.............#",
		"#.###########.#",
		"###############",
	}

	level, err := New(layout, 16)
	if err != nil {
		panic(err)
	}
	return level
}

// New builds a level from a slice of string rows and the tile size (pixels).
func New(layout []string, tileSize int) (*Level, error) {
	if len(layout) == 0 {
		return nil, fmt.Errorf("empty level layout")
	}

	height := len(layout)
	width := len(layout[0])
	tiles := make([][]TileType, height)
	warpEntrances := []GridPos{}

	for rowIdx, row := range layout {
		if len(row) != width {
			return nil, fmt.Errorf("inconsistent row width at %d", rowIdx)
		}

		tiles[rowIdx] = make([]TileType, width)
		for colIdx, ch := range row {
			switch ch {
			case '#':
				tiles[rowIdx][colIdx] = TileWall
			case '.':
				tiles[rowIdx][colIdx] = TilePath
			case 'W':
				tiles[rowIdx][colIdx] = TileWarp
				warpEntrances = append(warpEntrances, GridPos{Col: colIdx, Row: rowIdx})
			default:
				return nil, fmt.Errorf("unknown tile rune %q at row %d col %d", ch, rowIdx, colIdx)
			}
		}
	}

	warpTargets := map[GridPos]GridPos{}
	if len(warpEntrances)%2 != 0 {
		return nil, fmt.Errorf("warp entrances must be even, got %d", len(warpEntrances))
	}
	for i := 0; i < len(warpEntrances); i += 2 {
		a := warpEntrances[i]
		b := warpEntrances[i+1]
		warpTargets[a] = b
		warpTargets[b] = a
	}

	return &Level{
		Tiles:       tiles,
		TileSize:    tileSize,
		Width:       width,
		Height:      height,
		warpTargets: warpTargets,
	}, nil
}

// TileAt returns the tile type at the given column/row.
func (l *Level) TileAt(col, row int) TileType {
	if row < 0 || row >= l.Height || col < 0 || col >= l.Width {
		return TileWall
	}
	return l.Tiles[row][col]
}

// PixelWidth returns the stage width in pixels.
func (l *Level) PixelWidth() int {
	return l.Width * l.TileSize
}

// PixelHeight returns the stage height in pixels.
func (l *Level) PixelHeight() int {
	return l.Height * l.TileSize
}

// GridForPixel converts a pixel coordinate to a grid position.
func (l *Level) GridForPixel(x, y float64) GridPos {
	tile := float64(l.TileSize)
	return GridPos{
		Col: int(math.Floor(x / tile)),
		Row: int(math.Floor(y / tile)),
	}
}

// WarpTarget returns the paired warp destination for the provided grid cell.
func (l *Level) WarpTarget(pos GridPos) (GridPos, bool) {
	target, ok := l.warpTargets[pos]
	return target, ok
}

// Collides reports if a rectangle positioned at (x, y) with the given size is overlapping walls.
func (l *Level) Collides(x, y, size float64) bool {
	// Small epsilon reduces floating-point jitter at tile edges.
	const epsilon = 0.01
	points := [][2]float64{
		{x + epsilon, y + epsilon},
		{x + size - epsilon, y + epsilon},
		{x + epsilon, y + size - epsilon},
		{x + size - epsilon, y + size - epsilon},
	}

	for _, pt := range points {
		grid := l.GridForPixel(pt[0], pt[1])
		if l.TileAt(grid.Col, grid.Row) == TileWall {
			return true
		}
	}

	return false
}
