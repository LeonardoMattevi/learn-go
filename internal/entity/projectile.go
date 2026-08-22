package entity

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/LeonardoMattevi/go-game/internal/camera"
	"github.com/LeonardoMattevi/go-game/internal/world"
)

const projectileHalfSize = 3

type Projectile struct {
	X, Y   float64
	VX, VY float64
	Owner  OwnerType
	Damage int
	alive  bool
	img    *ebiten.Image
}

func NewProjectile(x, y, vx, vy float64, owner OwnerType, damage int, img *ebiten.Image) *Projectile {
	return &Projectile{X: x, Y: y, VX: vx, VY: vy, Owner: owner, Damage: damage, alive: true, img: img}
}

func (p *Projectile) IsAlive() bool { return p.alive }

func (p *Projectile) Kill() { p.alive = false }

func (p *Projectile) Bounds() image.Rectangle {
	x, y := int(p.X), int(p.Y)
	return image.Rect(x-projectileHalfSize, y-projectileHalfSize, x+projectileHalfSize, y+projectileHalfSize)
}

func (p *Projectile) Update(dt float64, w *world.World) error {
	p.X += p.VX * dt
	p.Y += p.VY * dt
	tx, ty := w.PixelToTile(p.X, p.Y)
	if w.IsSolid(tx, ty) || p.X < 0 || p.Y < 0 ||
		p.X > float64(w.PixelWidth()) || p.Y > float64(w.PixelHeight()) {
		p.alive = false
	}
	return nil
}

func (p *Projectile) Draw(screen *ebiten.Image, cam camera.Camera) {
	sx, sy := cam.WorldToScreen(p.X, p.Y)
	w, h := p.img.Bounds().Dx(), p.img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(p.img, op)
}

// RemoveDead filters out dead projectiles in-place without allocating a new slice.
func RemoveDead(ps []*Projectile) []*Projectile {
	n := 0
	for _, p := range ps {
		if p.alive {
			ps[n] = p
			n++
		}
	}
	return ps[:n]
}
