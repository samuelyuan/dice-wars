package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/samuelyuan/dice-wars/internal/app"
)

func main() {
	ebiten.SetWindowSize(app.DefaultScreenWidth, app.DefaultScreenHeight)
	ebiten.SetWindowTitle("Dice Wars")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(app.NewApp()); err != nil {
		log.Fatal(err)
	}
}
