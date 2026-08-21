package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const TileSize = 16

var (
	colorFloor = color.RGBA{R: 25, G: 25, B: 35, A: 255}
	colorWall  = color.RGBA{R: 70, G: 70, B: 85, A: 255}
)

type World struct {
	Tiles         [][]Tile
	Width, Height int // in tiles
}

func New(tiles [][]Tile) *World {
	return &World{
		Tiles:  tiles,
		Width:  len(tiles[0]),
		Height: len(tiles),
	}
}

func (w *World) PixelWidth() int  { return w.Width * TileSize }
func (w *World) PixelHeight() int { return w.Height * TileSize }

// IsSolid reports whether tile (tx, ty) blocks movement.
func (w *World) IsSolid(tx, ty int) bool {
	if tx < 0 || ty < 0 || tx >= w.Width || ty >= w.Height {
		return true
	}
	return w.Tiles[ty][tx] == TileWall
}

// PixelToTile converts pixel coords to tile coords.
func (w *World) PixelToTile(px, py float64) (int, int) {
	return int(px) / TileSize, int(py) / TileSize
}

// Draw renders the visible portion of the world.
// offsetX, offsetY are the camera's top-left in world-pixel space.
func (w *World) Draw(screen *ebiten.Image, offsetX, offsetY float64, screenW, screenH int) {
	for ty, row := range w.Tiles {
		for tx, tile := range row {
			sx := float64(tx*TileSize) - offsetX
			sy := float64(ty*TileSize) - offsetY
			if sx+TileSize < 0 || sy+TileSize < 0 || sx > float64(screenW) || sy > float64(screenH) {
				continue
			}
			var c color.RGBA
			switch tile {
			case TileFloor:
				c = colorFloor
			case TileWall:
				c = colorWall
			}
			vector.DrawFilledRect(screen, float32(sx), float32(sy), TileSize, TileSize, c, false)
		}
	}
}
