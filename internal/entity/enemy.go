package entity

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/LeonardoMattevi/go-game/internal/camera"
	"github.com/LeonardoMattevi/go-game/internal/world"
)

const (
	enemyHalfSize   = 6
	enemySpeed      = 55.0  // pixels/second
	walkerChaseRange  = 160.0
	shooterRange    = 180.0
	enemyShootDelay = 2.0 // seconds between shots
	enemyProjSpeed  = 110.0
)

type EnemyKind int

const (
	EnemyWalker  EnemyKind = iota
	EnemyShooter
)

type Enemy struct {
	X, Y          float64
	HP, MaxHP     int
	Dir           Direction
	Kind          EnemyKind
	shootCooldown float64
	frame         int
	frameTimer    float64
	alive         bool
}

func NewEnemy(x, y float64, kind EnemyKind) *Enemy {
	return &Enemy{
		X:     x,
		Y:     y,
		HP:    3,
		MaxHP: 3,
		Kind:  kind,
		Dir:   DirDown,
		alive: true,
	}
}

func (e *Enemy) IsAlive() bool { return e.alive && e.HP > 0 }

func (e *Enemy) Bounds() image.Rectangle {
	x, y := int(e.X), int(e.Y)
	return image.Rect(x-enemyHalfSize, y-enemyHalfSize, x+enemyHalfSize, y+enemyHalfSize)
}

// Update advances the enemy AI toward the player at (px, py).
// Returns a new Projectile if a Shooter fires this tick, otherwise nil.
// Direct-chase movement (BFS deferred to a later refactor).
func (e *Enemy) Update(dt float64, w *world.World, px, py float64) *Projectile {
	if !e.IsAlive() {
		return nil
	}

	dx := px - e.X
	dy := py - e.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	switch e.Kind {
	case EnemyWalker:
		if dist < walkerChaseRange && dist > 0 {
			nx, ny := dx/dist, dy/dist
			if newX := e.X + nx*enemySpeed*dt; !e.collidesAt(newX, e.Y, w) {
				e.X = newX
			}
			if newY := e.Y + ny*enemySpeed*dt; !e.collidesAt(e.X, newY, w) {
				e.Y = newY
			}
			e.updateDir(nx, ny)
			e.frameTimer += dt
			if e.frameTimer >= frameInterval {
				e.frameTimer = 0
				e.frame = (e.frame + 1) % framesPerDir
			}
		}

	case EnemyShooter:
		if e.shootCooldown > 0 {
			e.shootCooldown -= dt
		}
		if dist < shooterRange && dist > 0 && e.shootCooldown <= 0 {
			nx, ny := dx/dist, dy/dist
			e.updateDir(nx, ny)
			e.shootCooldown = enemyShootDelay
			return NewProjectile(e.X, e.Y, nx*enemyProjSpeed, ny*enemyProjSpeed, OwnerEnemy, 1)
		}
	}

	return nil
}

func (e *Enemy) TakeDamage(dmg int) {
	if !e.alive {
		return
	}
	e.HP -= dmg
	if e.HP <= 0 {
		e.alive = false
	}
}

func (e *Enemy) Draw(screen *ebiten.Image, cam camera.Camera) {
	if !e.alive {
		return
	}
	sx, sy := cam.WorldToScreen(e.X, e.Y)
	var col color.RGBA
	switch e.Kind {
	case EnemyWalker:
		col = color.RGBA{R: 220, G: 50, B: 50, A: 255}
	case EnemyShooter:
		col = color.RGBA{R: 200, G: 110, B: 20, A: 255}
	}
	vector.DrawFilledRect(
		screen,
		float32(sx)-enemyHalfSize, float32(sy)-enemyHalfSize,
		enemyHalfSize*2, enemyHalfSize*2,
		col,
		false,
	)
}

func (e *Enemy) collidesAt(x, y float64, w *world.World) bool {
	corners := [4][2]float64{
		{x - enemyHalfSize, y - enemyHalfSize},
		{x + enemyHalfSize, y - enemyHalfSize},
		{x - enemyHalfSize, y + enemyHalfSize},
		{x + enemyHalfSize, y + enemyHalfSize},
	}
	for _, c := range corners {
		tx, ty := w.PixelToTile(c[0], c[1])
		if w.IsSolid(tx, ty) {
			return true
		}
	}
	return false
}

func (e *Enemy) updateDir(nx, ny float64) {
	if math.Abs(ny) >= math.Abs(nx) {
		if ny < 0 {
			e.Dir = DirUp
		} else {
			e.Dir = DirDown
		}
	} else {
		if nx < 0 {
			e.Dir = DirLeft
		} else {
			e.Dir = DirRight
		}
	}
}

// RemoveDeadEnemies filters out dead enemies in-place without allocating.
func RemoveDeadEnemies(es []*Enemy) []*Enemy {
	n := 0
	for _, e := range es {
		if e.IsAlive() {
			es[n] = e
			n++
		}
	}
	return es[:n]
}
