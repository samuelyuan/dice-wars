package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/samuelyuan/dice-wars/internal/app"
)

func main() {
	ebiten.SetWindowSize(app.ScreenWidth, app.ScreenHeight)
	ebiten.SetWindowTitle("Dice Wars")

	if err := ebiten.RunGame(app.NewApp()); err != nil {
		log.Fatal(err)
	}
}
