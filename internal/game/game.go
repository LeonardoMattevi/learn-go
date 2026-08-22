package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/LeonardoMattevi/go-game/internal/camera"
	"github.com/LeonardoMattevi/go-game/internal/entity"
	"github.com/LeonardoMattevi/go-game/internal/sprite"
	"github.com/LeonardoMattevi/go-game/internal/world"
)

const (
	ScreenWidth      = 320
	ScreenHeight     = 240
	dt               = 1.0 / 60.0
	projectileSpeed  = 200.0
	shootCooldownMax = 0.25
)

type GameState int

const (
	StatePlaying       GameState = iota
	StateGameOver
	StateLevelComplete
)

// 40×30 tile map (640×480 pixel world, 320×240 visible = scrollable).
// 0 = floor, 1 = wall.
var levelTiles = [][]world.Tile{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

// tileCenter converts a tile column/row to world-pixel center coordinates.
func tileCenter(col, row int) (float64, float64) {
	return float64(col*world.TileSize + world.TileSize/2),
		float64(row*world.TileSize + world.TileSize/2)
}

func walker(col, row int, spr *sprite.Sprite, projImg *ebiten.Image) *entity.Enemy {
	x, y := tileCenter(col, row)
	return entity.NewEnemy(x, y, entity.EnemyWalker, spr, projImg)
}

func shooter(col, row int, spr *sprite.Sprite, projImg *ebiten.Image) *entity.Enemy {
	x, y := tileCenter(col, row)
	return entity.NewEnemy(x, y, entity.EnemyShooter, spr, projImg)
}

type Game struct {
	world         *world.World
	cam           camera.Camera
	player        *entity.Player
	projectiles   []*entity.Projectile
	enemies       []*entity.Enemy
	shootCooldown float64
	state         GameState
	imgProjPlayer *ebiten.Image
	imgProjEnemy  *ebiten.Image
	leftStick     *VirtualStick
	rightStick    *VirtualStick
}

func New() *Game {
	// Build all sprites once at startup.
	sprPlayer  := sprite.BuildPlayer()
	sprWalker  := sprite.BuildWalker()
	sprShooter := sprite.BuildShooter()
	imgProjPlayer := sprite.BuildProjectilePlayer()
	imgProjEnemy  := sprite.BuildProjectileEnemy()
	tileFloor := sprite.BuildTileFloor()
	tileWall  := sprite.BuildTileWall()

	px, py := tileCenter(2, 2)
	g := &Game{
		world:      world.New(levelTiles, tileFloor, tileWall),
		cam:        camera.Camera{ScreenW: ScreenWidth, ScreenH: ScreenHeight},
		player:     entity.NewPlayer(px, py, sprPlayer),
		leftStick:  newStick(false, 55, ScreenHeight-55),
		rightStick: newStick(true, ScreenWidth-55, ScreenHeight-55),
		enemies: []*entity.Enemy{
			walker(30, 2, sprWalker, imgProjEnemy),
			walker(30, 10, sprWalker, imgProjEnemy),
			walker(5, 22, sprWalker, imgProjEnemy),
			walker(20, 25, sprWalker, imgProjEnemy),
			shooter(15, 13, sprShooter, imgProjEnemy),
			shooter(32, 20, sprShooter, imgProjEnemy),
		},
		imgProjPlayer: imgProjPlayer,
		imgProjEnemy:  imgProjEnemy,
	}
	g.cam.Follow(g.player.X, g.player.Y, g.world.PixelWidth(), g.world.PixelHeight())
	return g
}

func (g *Game) Update() error {
	// Handle restart from non-playing states.
	if g.state != StatePlaying {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			*g = *New()
		}
		return nil
	}

	// Update joysticks and feed their output to the player.
	g.leftStick.Update(ScreenWidth)
	g.rightStick.Update(ScreenWidth)
	g.player.JoyMoveDX = g.leftStick.DX()
	g.player.JoyMoveDY = g.leftStick.DY()
	g.player.JoyAimDX = g.rightStick.DX()
	g.player.JoyAimDY = g.rightStick.DY()
	g.player.JoyAimActive = g.rightStick.Active()

	if err := g.player.Update(dt, g.world); err != nil {
		return err
	}

	// Shoot cooldown
	if g.shootCooldown > 0 {
		g.shootCooldown -= dt
	}

	// Auto-fire: arrow keys (desktop) or right joystick (touch).
	arrowHeld := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	if (arrowHeld || g.rightStick.Active()) && g.shootCooldown <= 0 {
		g.projectiles = append(g.projectiles, entity.NewProjectile(
			g.player.X, g.player.Y,
			g.player.ShootDX*projectileSpeed,
			g.player.ShootDY*projectileSpeed,
			entity.OwnerPlayer, 1, g.imgProjPlayer,
		))
		g.shootCooldown = shootCooldownMax
	}

	// Update enemies; collect any projectiles they fire.
	for _, e := range g.enemies {
		if proj := e.Update(dt, g.world, g.player.X, g.player.Y); proj != nil {
			g.projectiles = append(g.projectiles, proj)
		}
	}

	// Update projectiles.
	for _, p := range g.projectiles {
		p.Update(dt, g.world)
	}

	// Collision: projectiles vs enemies and player.
	playerBounds := g.player.Bounds()
	for _, proj := range g.projectiles {
		if !proj.IsAlive() {
			continue
		}
		pb := proj.Bounds()
		if proj.Owner == entity.OwnerPlayer {
			for _, e := range g.enemies {
				if e.IsAlive() && pb.Overlaps(e.Bounds()) {
					e.TakeDamage(proj.Damage)
					proj.Kill()
					break
				}
			}
		} else {
			if pb.Overlaps(playerBounds) {
				g.player.TakeDamage(proj.Damage)
				proj.Kill()
			}
		}
	}

	// Walker contact damage to player.
	for _, e := range g.enemies {
		if e.IsAlive() && e.Kind == entity.EnemyWalker && e.Bounds().Overlaps(playerBounds) {
			g.player.TakeDamage(1)
		}
	}

	g.projectiles = entity.RemoveDead(g.projectiles)
	g.enemies = entity.RemoveDeadEnemies(g.enemies)

	// State transitions.
	if g.player.HP <= 0 {
		g.state = StateGameOver
	} else if len(g.enemies) == 0 {
		g.state = StateLevelComplete
	}

	g.cam.Follow(g.player.X, g.player.Y, g.world.PixelWidth(), g.world.PixelHeight())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 10, G: 10, B: 20, A: 255})
	g.world.Draw(screen, g.cam.X, g.cam.Y, ScreenWidth, ScreenHeight)
	for _, p := range g.projectiles {
		p.Draw(screen, g.cam)
	}
	for _, e := range g.enemies {
		e.Draw(screen, g.cam)
	}
	g.player.Draw(screen, g.cam)
	g.drawHUD(screen)
	g.leftStick.Draw(screen)
	g.rightStick.Draw(screen)
	if g.state != StatePlaying {
		g.drawOverlay(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
