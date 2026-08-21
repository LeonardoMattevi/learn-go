package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/LeonardoMattevi/go-game/internal/game"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Druid II")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
