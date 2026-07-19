package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/samuelyuan/dice-wars/internal/game"
)

const (
	mapDiceSize         = 28.0
	mapDiceHeightFactor = 0.54
	mapDiceWidthFactor  = 0.55
	mapDiceShiftX       = -5.0
	battleDiceSize      = 28.0
	battleDiceSpacing   = 4.0
	battleFirstSpacing  = 8.0
	battlePanelWidth    = 520.0
)

func drawDiceRoll(screen *ebiten.Image, board *game.Board, lc *LayoutContext) {
	attack := board.LastAttack
	layout := lc.computeBattlePanelLayout()
	centerX := float64(lc.Width / 2)
	centerY := layout.centerY

	drawBattlePanel(screen, centerX, centerY, layout.height)

	attackerRevealed := board.RevealedAttackerDice()
	defenderRevealed := board.RevealedDefenderDice()
	attackerOwner := board.AttackAttackerOwner()
	defenderOwner := board.AttackDefenderOwner()

	if attackerOwner >= 0 && attackerRevealed > 0 {
		drawBattleSide(screen, centerX, centerY, layout.diceSize, attack.AttackerRolls, attackerRevealed, attack.AttackTotal, attackerOwner, true)
	}
	if defenderOwner >= 0 && defenderRevealed > 0 {
		drawBattleSide(screen, centerX, centerY, layout.diceSize, attack.DefenderRolls, defenderRevealed, attack.DefenseTotal, defenderOwner, false)
	}
	if attackerRevealed == 0 && defenderRevealed == 0 {
		textWidth := len("Rolling...") * 7
		drawText(screen, "Rolling...", int(centerX)-textWidth/2, int(centerY)+4, colorText)
	}
}

func drawBattlePanel(screen *ebiten.Image, centerX, centerY, height float64) {
	panelW := float32(battlePanelWidth)
	panelH := float32(height)
	left := float32(centerX) - panelW/2
	top := float32(centerY) - panelH/2
	vector.DrawFilledRect(screen, left, top, panelW, panelH, color.RGBA{245, 245, 245, 240}, false)
	vector.StrokeRect(screen, left, top, panelW, panelH, 2, color.Black, false)
}

func drawBattleSide(screen *ebiten.Image, centerX, centerY, diceSize float64, rolls []int, revealed, total, playerIdx int, leftSide bool) {
	dieY := centerY - diceSize/2
	step := diceSize + battleDiceSpacing
	rollCount := len(rolls)

	for i := 0; i < revealed && i < rollCount; i++ {
		slot := rollCount - 1 - i
		dieX := battleDieX(centerX, step, slot, leftSide)
		drawPlayerDie(screen, dieX, dieY, diceSize, rolls[slot], playerIdx)
	}

	if revealed != rollCount {
		return
	}
	scoreY := textCenterY(int(dieY), int(diceSize))
	scoreX := battleScoreX(centerX, step, rollCount, leftSide)
	drawText(screen, strconv.Itoa(total), scoreX, scoreY, colorText)
}

func battleDieX(centerX, step float64, slot int, leftSide bool) float64 {
	if leftSide {
		return centerX - battleFirstSpacing - step*float64(slot+1)
	}
	return centerX + battleFirstSpacing + step*float64(slot)
}

func battleScoreX(centerX, step float64, rollCount int, leftSide bool) int {
	if leftSide {
		return int(centerX - battleFirstSpacing - step*float64(rollCount+1))
	}
	return int(centerX + battleFirstSpacing + step*float64(rollCount) + 4)
}
