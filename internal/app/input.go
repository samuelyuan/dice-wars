package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/samuelyuan/dice-wars/internal/game"
)

func (a *App) handleGameInput() {
	mx, my := ebiten.CursorPosition()
	a.hoverBtn = a.hoveredButton(mx, my)

	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	switch a.hoverBtn {
	case "menu":
		a.wantMenu = true
	case "end":
		a.board.EndTurn()
	case "auto":
		a.board.StartAITurn()
	case "cheat":
		a.board.CheatMode = !a.board.CheatMode
	default:
		if my < int(mapContentBottom()) && !a.board.IsBusy() {
			a.board.Click(float64(mx)-MapOffsetX, float64(my)-MapOffsetY)
		}
	}
}

func (a *App) hoveredButton(mx, my int) string {
	if BtnMenu.Contains(mx, my) {
		return "menu"
	}
	if a.board.IsHumanTurn() {
		if BtnEndTurn.Contains(mx, my) {
			return "end"
		}
		if BtnAuto.Contains(mx, my) {
			return "auto"
		}
	}
	if CheatHit.Contains(mx, my) {
		return "cheat"
	}
	return ""
}

func drawGame(screen *ebiten.Image, board *game.Board, hoverBtn string) {
	drawMap(screen, board)
	drawGameHUD(screen, board, hoverBtn)
}
