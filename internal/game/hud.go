package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	hudBarW   = 72
	hudBarH   = 6
	hudMargin = 4
)

func (g *Game) drawHUD(screen *ebiten.Image) {
	// HP bar
	if g.player.MaxHP > 0 {
		hpRatio := float32(g.player.HP) / float32(g.player.MaxHP)
		vector.DrawFilledRect(screen, hudMargin, hudMargin, hudBarW, hudBarH,
			color.RGBA{60, 20, 20, 220}, false)
		vector.DrawFilledRect(screen, hudMargin, hudMargin, float32(hudBarW)*hpRatio, hudBarH,
			color.RGBA{220, 50, 50, 255}, false)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %d/%d", g.player.HP, g.player.MaxHP),
		hudMargin, hudMargin+hudBarH+1)

	// Mana bar
	const manaY = hudMargin + hudBarH + 14
	if g.player.MaxMana > 0 {
		mpRatio := float32(g.player.Mana) / float32(g.player.MaxMana)
		vector.DrawFilledRect(screen, hudMargin, manaY, hudBarW, hudBarH,
			color.RGBA{20, 20, 60, 220}, false)
		vector.DrawFilledRect(screen, hudMargin, manaY, float32(hudBarW)*mpRatio, hudBarH,
			color.RGBA{60, 110, 220, 255}, false)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MP %d/%d", g.player.Mana, g.player.MaxMana),
		hudMargin, manaY+hudBarH+1)

	// Enemy count (top-right)
	label := fmt.Sprintf("Enemies: %d", len(g.enemies))
	ebitenutil.DebugPrintAt(screen, label, ScreenWidth-len(label)*6-hudMargin, hudMargin)
}

func (g *Game) drawOverlay(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight,
		color.RGBA{0, 0, 0, 160}, false)

	cy := ScreenHeight/2 - 8
	switch g.state {
	case StateGameOver:
		printCentered(screen, "GAME OVER", cy-8)
		printCentered(screen, "Tap or press R to restart", cy+8)
	case StateLevelComplete:
		printCentered(screen, "LEVEL COMPLETE!", cy-8)
		printCentered(screen, "Tap or press R to restart", cy+8)
	}
}

// printCentered draws text horizontally centered on screen using the debug font (6px/char).
func printCentered(screen *ebiten.Image, s string, y int) {
	x := (ScreenWidth - len(s)*6) / 2
	ebitenutil.DebugPrintAt(screen, s, x, y)
}
