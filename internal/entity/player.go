package entity

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/LeonardoMattevi/go-game/internal/camera"
	"github.com/LeonardoMattevi/go-game/internal/sprite"
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
	ShootDX       float64
	ShootDY       float64
	Invincible    float64
	// Joystick input set by game.Update before calling Update().
	JoyMoveDX, JoyMoveDY float64
	JoyAimDX, JoyAimDY   float64
	JoyAimActive         bool
	frame      int
	frameTimer float64
	spr        *sprite.Sprite
}

func NewPlayer(x, y float64, spr *sprite.Sprite) *Player {
	return &Player{
		X:       x,
		Y:       y,
		HP:      10,
		MaxHP:   10,
		Mana:    10,
		MaxMana: 10,
		Dir:     DirDown,
		ShootDY: 1,
		spr:     spr,
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

	// Movement: WASD only.
	up := ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyD)

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

	// Joystick overrides keyboard when left stick is active.
	if p.JoyMoveDX != 0 || p.JoyMoveDY != 0 {
		dx = p.JoyMoveDX
		dy = p.JoyMoveDY
		// Dir from joystick angle (compare squares to avoid math.Abs import).
		if dy*dy >= dx*dx {
			if dy < 0 {
				p.Dir = DirUp
			} else {
				p.Dir = DirDown
			}
		} else {
			if dx < 0 {
				p.Dir = DirLeft
			} else {
				p.Dir = DirRight
			}
		}
	}

	// Normalize diagonal so speed is consistent in all directions.
	if dx != 0 && dy != 0 {
		dx *= 0.7071
		dy *= 0.7071
	}

	// Aim direction: arrow keys, independent of movement.
	// Retains last aimed direction when no arrow is held.
	aL := ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	aR := ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	aU := ebiten.IsKeyPressed(ebiten.KeyArrowUp)
	aD := ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	var sdx, sdy float64
	if aL && !aR {
		sdx = -1
	}
	if aR && !aL {
		sdx = 1
	}
	if aU && !aD {
		sdy = -1
	}
	if aD && !aU {
		sdy = 1
	}
	if sdx != 0 && sdy != 0 {
		sdx *= 0.7071
		sdy *= 0.7071
	}
	if sdx != 0 || sdy != 0 {
		p.ShootDX = sdx
		p.ShootDY = sdy
	}
	// Right joystick overrides arrow-key aim when active.
	if p.JoyAimActive {
		p.ShootDX = p.JoyAimDX
		p.ShootDY = p.JoyAimDY
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
	// Flash during invincibility frames (hide every other 0.1s block).
	if p.Invincible > 0 && int(p.Invincible*10)%2 == 0 {
		return
	}
	p.spr.Draw(screen, sx, sy, int(p.Dir), p.frame)
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
