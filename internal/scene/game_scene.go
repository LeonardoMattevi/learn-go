package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/LeonardoMattevi/go-game/internal/game"
)

// GameScene wraps game.Game to satisfy the Scene interface.
// All gameplay logic (entities, HUD, joysticks, state machine) lives in game.Game.
type GameScene struct {
	g *game.Game
}

func NewGameScene() *GameScene {
	return &GameScene{g: game.New()}
}

func (s *GameScene) Update() (Scene, error) {
	return s, s.g.Update()
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	s.g.Draw(screen)
}
