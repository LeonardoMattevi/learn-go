package entity

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/LeonardoMattevi/go-game/internal/camera"
	"github.com/LeonardoMattevi/go-game/internal/world"
)

const (
	playerHalfSize  = 6   // half of the 12×12 bounding box
	playerSpeed     = 120 // pixels per second
	frameInterval   = 0.15
	framesPerDir    = 4
)

type Player struct {
	X, Y          float64
	HP, MaxHP     int
	Mana, MaxMana int
	Dir           Direction
	ShootDX       float64 // normalized shoot direction, updated when moving
	ShootDY       float64
	Invincible    float64 // seconds of invincibility remaining after taking damage
	frame         int
	frameTimer    float64
	// sprite *ebiten.Image  // wired in Fase 8
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		X:       x,
		Y:       y,
		HP:      10,
		MaxHP:   10,
		Mana:    10,
		MaxMana: 10,
		Dir:     DirDown,
		ShootDY: 1, // default facing down
	}
}

func (p *Player) IsAlive() bool { return p.HP > 0 }

func (p *Player) TakeDamage(dmg int) {
	if p.Invincible > 0 {
		return
	}
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
	p.Invincible = 1.0
}

func (p *Player) Bounds() image.Rectangle {
	x, y := int(p.X), int(p.Y)
	return image.Rect(x-playerHalfSize, y-playerHalfSize, x+playerHalfSize, y+playerHalfSize)
}

func (p *Player) Update(dt float64, w *world.World) error {
	if p.Invincible > 0 {
		p.Invincible -= dt
	}

	up := ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)

	var dx, dy float64

	// Opposite keys cancel each other out.
	if left && !right {
		dx = -1
		p.Dir = DirLeft
	}
	if right && !left {
		dx = 1
		p.Dir = DirRight
	}
	// Vertical wins over horizontal for Dir when both axes are active.
	if up && !down {
		dy = -1
		p.Dir = DirUp
	}
	if down && !up {
		dy = 1
		p.Dir = DirDown
	}

	// Normalize diagonal so speed is consistent in all directions.
	if dx != 0 && dy != 0 {
		dx *= 0.7071
		dy *= 0.7071
	}

	// Track shoot direction from actual movement vector.
	if dx != 0 || dy != 0 {
		p.ShootDX = dx
		p.ShootDY = dy
	}

	if dx != 0 || dy != 0 {
		p.frameTimer += dt
		if p.frameTimer >= frameInterval {
			p.frameTimer = 0
			p.frame = (p.frame + 1) % framesPerDir
		}
	} else {
		p.frame = 0
		p.frameTimer = 0
	}

	// Move X, then Y separately so the player slides along walls.
	if newX := p.X + dx*playerSpeed*dt; !p.collidesAt(newX, p.Y, w) {
		p.X = newX
	}
	if newY := p.Y + dy*playerSpeed*dt; !p.collidesAt(p.X, newY, w) {
		p.Y = newY
	}

	return nil
}

func (p *Player) Draw(screen *ebiten.Image, cam camera.Camera) {
	sx, sy := cam.WorldToScreen(p.X, p.Y)
	vector.DrawFilledRect(
		screen,
		float32(sx)-playerHalfSize, float32(sy)-playerHalfSize,
		playerHalfSize*2, playerHalfSize*2,
		color.RGBA{R: 0, G: 200, B: 80, A: 255},
		false,
	)
}

// collidesAt checks whether the player's bounding box at (x, y) overlaps a solid tile.
func (p *Player) collidesAt(x, y float64, w *world.World) bool {
	corners := [4][2]float64{
		{x - playerHalfSize, y - playerHalfSize},
		{x + playerHalfSize, y - playerHalfSize},
		{x - playerHalfSize, y + playerHalfSize},
		{x + playerHalfSize, y + playerHalfSize},
	}
	for _, c := range corners {
		tx, ty := w.PixelToTile(c[0], c[1])
		if w.IsSolid(tx, ty) {
			return true
		}
	}
	return false
}
