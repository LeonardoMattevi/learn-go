package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const TileSize = 16

type World struct {
	Tiles         [][]Tile
	Width, Height int
	tileFloor     *ebiten.Image
	tileWall      *ebiten.Image
}

func New(tiles [][]Tile, floor, wall *ebiten.Image) *World {
	return &World{
		Tiles:     tiles,
		Width:     len(tiles[0]),
		Height:    len(tiles),
		tileFloor: floor,
		tileWall:  wall,
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

// Draw renders the visible portion of the world using tile images.
// offsetX, offsetY are the camera's top-left in world-pixel space.
func (w *World) Draw(screen *ebiten.Image, offsetX, offsetY float64, screenW, screenH int) {
	for ty, row := range w.Tiles {
		for tx, tile := range row {
			sx := float64(tx*TileSize) - offsetX
			sy := float64(ty*TileSize) - offsetY
			if sx+TileSize < 0 || sy+TileSize < 0 || sx > float64(screenW) || sy > float64(screenH) {
				continue
			}
			var img *ebiten.Image
			switch tile {
			case TileFloor:
				img = w.tileFloor
			case TileWall:
				img = w.tileWall
			}
			if img == nil {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(sx, sy)
			screen.DrawImage(img, op)
		}
	}
}
