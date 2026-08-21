package entity

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/LeonardoMattevi/go-game/internal/camera"
	"github.com/LeonardoMattevi/go-game/internal/world"
)

// Entity is implemented by everything that lives in the world:
// Player, Enemy, Projectile.
type Entity interface {
	Update(dt float64, w *world.World) error
	Draw(screen *ebiten.Image, cam camera.Camera)
	Bounds() image.Rectangle // bounding box in world-pixel coords
	IsAlive() bool
}

type Direction int

const (
	DirDown Direction = iota
	DirUp
	DirLeft
	DirRight
)

type OwnerType int

const (
	OwnerPlayer OwnerType = iota
	OwnerEnemy
)
