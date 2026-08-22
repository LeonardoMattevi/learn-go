package scene

import "github.com/hajimehoshi/ebiten/v2"

// Scene is any full-screen game state that can update and draw itself.
// Update returns the scene to run next frame — return self to stay, return
// a new scene to transition immediately.
type Scene interface {
	Update() (Scene, error)
	Draw(screen *ebiten.Image)
}

// Root implements ebiten.Game and delegates to the active Scene.
type Root struct {
	current Scene
}

func NewRoot() *Root {
	return &Root{current: NewMenu()}
}

func (r *Root) Update() error {
	next, err := r.current.Update()
	if next != nil {
		r.current = next
	}
	return err
}

func (r *Root) Draw(screen *ebiten.Image) {
	r.current.Draw(screen)
}

func (r *Root) Layout(outsideW, outsideH int) (int, int) {
	return 320, 240
}
