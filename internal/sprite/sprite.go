package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite holds a set of animation frames and draws the current one.
// Frames are grouped by direction: frames[dir*FramesPerDir + frame].
// If the sprite has no directional animation (single group), use dir=0.
type Sprite struct {
	frames      []*ebiten.Image
	dirsCount   int // number of direction groups
	framesPerDir int
}

// New creates a Sprite from a flat slice of frames.
// dirsCount=1 means no directional animation.
func New(frames []*ebiten.Image, dirsCount, framesPerDir int) *Sprite {
	return &Sprite{frames: frames, dirsCount: dirsCount, framesPerDir: framesPerDir}
}

// Draw renders the sprite centered on (sx, sy) in screen space.
// dir selects the direction group; frame selects the animation frame within it.
func (s *Sprite) Draw(screen *ebiten.Image, sx, sy float64, dir, frame int) {
	if len(s.frames) == 0 {
		return
	}
	idx := dir*s.framesPerDir + frame
	if idx < 0 || idx >= len(s.frames) {
		idx = 0
	}
	img := s.frames[idx]
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(img, op)
}
