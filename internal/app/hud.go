package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/samuelyuan/dice-wars/internal/game"
)

func drawGameHUD(screen *ebiten.Image, board *game.Board, hoverBtn string, lc *LayoutContext, fastForward bool) {
	if board.Phase == game.PhaseDiceRoll {
		drawDiceRoll(screen, board, lc)
	}

	drawTurnBanner(screen, board, lc)
	drawPlayerBar(screen, board, lc)

	if board.StatusMessage != "" {
		drawText(screen, board.StatusMessage, textCenterX(0, lc.Width, board.StatusMessage), int(lc.StatusTextY()), colorText)
	}

	if board.IsHumanTurn() {
		lc.BtnEndTurn().Draw(screen, hoverBtn == "end")
		lc.BtnAuto().Draw(screen, hoverBtn == "auto")
	}
	lc.BtnMenu().Draw(screen, hoverBtn == "menu")

	// Fast-forward stays visually "pressed" while active, not just on hover.
	lc.BtnFastForward().Draw(screen, fastForward || hoverBtn == "fastforward")

	if board.CheatMode {
		cheatRect := lc.CheatHit()
		drawText(screen, "Cheater!", cheatRect.X, cheatRect.Y+14, colorText)
	}
}

func drawTurnBanner(screen *ebiten.Image, board *game.Board, lc *LayoutContext) {
	banner := board.TurnBanner()
	player := board.Players[board.PlayerTurn]

	bx := float32(lc.TurnBannerX())
	by := float32(lc.TurnBannerY())
	vector.DrawFilledRect(screen, bx, by, TurnBannerW, TurnBannerH, color.RGBA{248, 248, 248, 255}, false)
	vector.StrokeRect(screen, bx, by, TurnBannerW, TurnBannerH, 2, color.Black, false)
	drawPlayerIcon(screen, float64(bx)+18, float64(by)+16, 20, player.Index)
	drawText(screen, banner, int(bx)+34, int(by)+22, colorText)
}

func drawPlayerBar(screen *ebiten.Image, board *game.Board, lc *LayoutContext) {
	startX := (lc.Width - board.NumPlayers*PlayerBarSlotW) / 2

	for i := 0; i < board.NumPlayers; i++ {
		if !board.IsPlayerActive(i) {
			continue
		}
		drawPlayerBarSlot(screen, board, i, startX+i*PlayerBarSlotW, lc)
	}
}

func drawPlayerBarSlot(screen *ebiten.Image, board *game.Board, playerIdx, x int, lc *LayoutContext) {
	bg := color.RGBA{255, 255, 255, 255}
	if playerIdx == board.PlayerTurn {
		bg = color.RGBA{255, 255, 0, 255}
		if board.Phase != game.PhaseIdle {
			bg = color.RGBA{255, 220, 0, 255}
		}
	}
	slotW := PlayerBarSlotW - 8
	playerBarY := int(lc.PlayerBarY())
	vector.DrawFilledRect(screen, float32(x), float32(playerBarY), float32(slotW), PlayerBarHeight, bg, false)
	vector.StrokeRect(screen, float32(x), float32(playerBarY), float32(slotW), PlayerBarHeight, 2, color.Black, false)
	drawPlayerIcon(screen, float64(x+22), float64(playerBarY+PlayerBarHeight/2), 28, playerIdx)
	drawText(screen, strconv.Itoa(board.ConnectedTerrCount(playerIdx)), x+42, playerBarY+26, colorText)
}
