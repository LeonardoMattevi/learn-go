package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW = 320
	screenH = 240
)

// btn is a rectangular button that responds to mouse hover, keyboard, and touch.
type btn struct {
	label    string
	x, y     int
	w, h     int
}

func (b *btn) contains(px, py int) bool {
	return px >= b.x && px < b.x+b.w && py >= b.y && py < b.y+b.h
}

func (b *btn) draw(screen *ebiten.Image, hover bool) {
	bg := color.RGBA{R: 50, G: 50, B: 80, A: 210}
	border := color.RGBA{R: 120, G: 120, B: 180, A: 255}
	if hover {
		bg = color.RGBA{R: 80, G: 80, B: 140, A: 230}
		border = color.RGBA{R: 180, G: 180, B: 255, A: 255}
	}
	// Background
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), bg, false)
	// Border (1px on each side via four thin rects)
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), 1, border, false)
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y+b.h-1), float32(b.w), 1, border, false)
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), 1, float32(b.h), border, false)
	vector.DrawFilledRect(screen, float32(b.x+b.w-1), float32(b.y), 1, float32(b.h), border, false)
	// Label centered
	tx := b.x + (b.w-len(b.label)*6)/2
	ty := b.y + (b.h-13)/2
	ebitenutil.DebugPrintAt(screen, b.label, tx, ty)
}

// MenuScene is the title screen shown at startup.
type MenuScene struct {
	startBtn      btn
	fullscreenBtn btn
}

func NewMenu() *MenuScene {
	cx := screenW / 2
	return &MenuScene{
		startBtn:      btn{"START GAME", cx - 55, 110, 110, 22},
		fullscreenBtn: btn{"FULLSCREEN", cx - 55, 145, 110, 22},
	}
}

func (m *MenuScene) Update() (Scene, error) {
	// Keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return NewGameScene(), nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Mouse click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if m.startBtn.contains(mx, my) {
			return NewGameScene(), nil
		}
		if m.fullscreenBtn.contains(mx, my) {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
	}

	// Touch tap
	var ids []ebiten.TouchID
	ids = inpututil.AppendJustPressedTouchIDs(ids)
	for _, id := range ids {
		tx, ty := ebiten.TouchPosition(id)
		if m.startBtn.contains(tx, ty) {
			return NewGameScene(), nil
		}
		if m.fullscreenBtn.contains(tx, ty) {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
	}

	return m, nil
}

func (m *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 10, G: 10, B: 20, A: 255})

	// Decorative header band
	vector.DrawFilledRect(screen, 0, 28, screenW, 2, color.RGBA{R: 80, G: 80, B: 140, A: 255}, false)
	vector.DrawFilledRect(screen, 0, 60, screenW, 2, color.RGBA{R: 80, G: 80, B: 140, A: 255}, false)

	// Title with shadow
	const title = "DRUID  II"
	tx := (screenW - len(title)*6) / 2
	ebitenutil.DebugPrintAt(screen, title, tx+1, 38+1) // shadow
	ebitenutil.DebugPrintAt(screen, title, tx, 38)

	// Subtitle
	const sub = "a Gauntlet-style adventure"
	ebitenutil.DebugPrintAt(screen, sub, (screenW-len(sub)*6)/2, 70)

	// Buttons with mouse-hover highlight
	mx, my := ebiten.CursorPosition()
	m.startBtn.draw(screen, m.startBtn.contains(mx, my))
	m.fullscreenBtn.draw(screen, m.fullscreenBtn.contains(mx, my))

	// Controls hint
	const hint = "WASD + arrows  /  dual stick on touch"
	ebitenutil.DebugPrintAt(screen, hint, (screenW-len(hint)*6)/2, 200)
}
