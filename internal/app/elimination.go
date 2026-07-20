package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/samuelyuan/dice-wars/internal/game"
)

// handleEliminationOverlayInput processes clicks on the overlay's two
// buttons: restart with the same player settings, or return to the menu.
func (a *App) handleEliminationOverlayInput() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	mx, my := ebiten.CursorPosition()
	lc := a.layout

	switch {
	case lc.EliminationRestartButton().Contains(mx, my):
		a.board = game.NewBoard(a.menu.NumPlayers, a.menu.HumanList())
		a.humanEliminated = false
		a.showEliminationOverlay = false
	case lc.EliminationNewGameButton().Contains(mx, my):
		a.board = nil
		a.humanEliminated = false
		a.showEliminationOverlay = false
		a.screen = ScreenMenu
	}
}

func (a *App) drawEliminationOverlay(screen *ebiten.Image, lc *LayoutContext) {
	vector.DrawFilledRect(screen, 0, 0, float32(lc.Width), float32(lc.Height), color.RGBA{0, 0, 0, 160}, false)

	msg := "You have been eliminated!"
	centerY := lc.Height/2 - 40
	drawText(screen, msg, textCenterX(0, lc.Width, msg), centerY, colorTextLight)

	sub := "Restart or return to the menu."
	drawText(screen, sub, textCenterX(0, lc.Width, sub), centerY+24, colorTextLight)

	mx, my := ebiten.CursorPosition()
	lc.EliminationRestartButton().Draw(screen, lc.EliminationRestartButton().Contains(mx, my))
	lc.EliminationNewGameButton().Draw(screen, lc.EliminationNewGameButton().Contains(mx, my))
}
