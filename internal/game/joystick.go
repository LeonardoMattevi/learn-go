package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	stickRadius   = float32(40)
	knobRadius    = float32(13)
	stickHintAlpha = uint8(40)  // ring opacity when not touched
	stickRingAlpha = uint8(90)  // ring opacity when active
	knobAlpha     = uint8(180)
)

// VirtualStick is a single on-screen analog stick driven by touch input.
// Each stick "owns" one half of the screen (left or right).
type VirtualStick struct {
	isRight  bool          // true = right half of screen
	hintX    float32       // fixed position shown as hint when not touched
	hintY    float32
	baseX    float32       // center of ring (set on touch-down)
	baseY    float32
	knobX    float32
	knobY    float32
	touchID  ebiten.TouchID
	active   bool
}

func newStick(isRight bool, hintX, hintY float32) *VirtualStick {
	return &VirtualStick{
		isRight: isRight,
		hintX:   hintX,
		hintY:   hintY,
		knobX:   hintX,
		knobY:   hintY,
		baseX:   hintX,
		baseY:   hintY,
	}
}

// Update reads touch events and updates the stick state.
func (s *VirtualStick) Update(screenW int) {
	var ids []ebiten.TouchID
	ids = ebiten.AppendTouchIDs(ids)

	// If we already own a touch, track it or release it.
	if s.active {
		found := false
		for _, id := range ids {
			if id == s.touchID {
				found = true
				tx, ty := ebiten.TouchPosition(id)
				s.moveKnob(float32(tx), float32(ty))
				break
			}
		}
		if !found {
			s.release()
		}
		return
	}

	// Claim the first unowned touch in our half of the screen.
	for _, id := range ids {
		tx, ty := ebiten.TouchPosition(id)
		inRight := tx > screenW/2
		if inRight == s.isRight {
			s.touchID = id
			s.baseX = float32(tx)
			s.baseY = float32(ty)
			s.knobX = float32(tx)
			s.knobY = float32(ty)
			s.active = true
			break
		}
	}
}

func (s *VirtualStick) release() {
	s.active = false
	s.baseX = s.hintX
	s.baseY = s.hintY
	s.knobX = s.hintX
	s.knobY = s.hintY
}

func (s *VirtualStick) moveKnob(tx, ty float32) {
	dx := tx - s.baseX
	dy := ty - s.baseY
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if dist > stickRadius {
		dx = dx / dist * stickRadius
		dy = dy / dist * stickRadius
	}
	s.knobX = s.baseX + dx
	s.knobY = s.baseY + dy
}

// DX returns the horizontal axis in [-1, 1].
func (s *VirtualStick) DX() float64 {
	if !s.active {
		return 0
	}
	return float64((s.knobX - s.baseX) / stickRadius)
}

// DY returns the vertical axis in [-1, 1].
func (s *VirtualStick) DY() float64 {
	if !s.active {
		return 0
	}
	return float64((s.knobY - s.baseY) / stickRadius)
}

// Active reports whether a finger is currently on this stick.
func (s *VirtualStick) Active() bool { return s.active }

// Draw renders the stick in screen space.
func (s *VirtualStick) Draw(screen *ebiten.Image) {
	cx, cy := s.baseX, s.baseY
	ringAlpha := stickHintAlpha
	if s.active {
		ringAlpha = stickRingAlpha
	}
	// Outer ring
	vector.StrokeCircle(screen, cx, cy, stickRadius, 2,
		color.RGBA{R: 255, G: 255, B: 255, A: ringAlpha}, true)
	// Knob (only when touched)
	if s.active {
		vector.DrawFilledCircle(screen, s.knobX, s.knobY, knobRadius,
			color.RGBA{R: 255, G: 255, B: 255, A: knobAlpha}, true)
	}
}
