package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/samuelyuan/dice-wars/internal/game"
)

func (a *App) handleGameInput() {
	mx, my := ebiten.CursorPosition()
	a.hoverBtn = a.hoveredButton(mx, my, a.layout)

	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	switch a.hoverBtn {
	case "menu":
		a.wantMenu = true
	case "fastforward":
		a.fastForward = !a.fastForward
	case "end":
		a.board.EndTurn()
	case "auto":
		a.board.StartAITurn()
	case "cheat":
		a.board.CheatMode = !a.board.CheatMode
	default:
		if my < int(a.layout.mapContentBottom()) && !a.board.IsBusy() {
			a.board.Click(float64(mx)-a.layout.MapOffsetX(), float64(my)-a.layout.MapOffsetY())
		}
	}
}

func (a *App) hoveredButton(mx, my int, lc *LayoutContext) string {
	if lc.BtnMenu().Contains(mx, my) {
		return "menu"
	}
	if lc.BtnFastForward().Contains(mx, my) {
		return "fastforward"
	}
	if a.board.IsHumanTurn() {
		if lc.BtnEndTurn().Contains(mx, my) {
			return "end"
		}
		if lc.BtnAuto().Contains(mx, my) {
			return "auto"
		}
	}
	if lc.CheatHit().Contains(mx, my) {
		return "cheat"
	}
	return ""
}

func drawGame(screen *ebiten.Image, board *game.Board, hoverBtn string, lc *LayoutContext, fastForward bool) {
	drawMap(screen, board, lc)
	drawGameHUD(screen, board, hoverBtn, lc, fastForward)
}
