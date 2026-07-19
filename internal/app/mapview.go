package app

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/samuelyuan/dice-wars/internal/game"
)

func drawMap(screen *ebiten.Image, board *game.Board, lc *LayoutContext) {
	offsetX, offsetY := lc.MapOffsetX(), lc.MapOffsetY()
	order := drawOrderTerritories(board.Territories, board)

	for _, t := range order {
		highlighted := territoryHighlighted(board, t)
		fill := playerFill(t.Owner, highlighted)
		for _, ax := range t.CellIDs {
			h := board.Grid.Hexes[ax]
			drawHexFill(screen, h, offsetX, offsetY, fill)
		}
	}

	for _, t := range order {
		for _, ax := range t.CellIDs {
			h := board.Grid.Hexes[ax]
			drawHexBorders(screen, board, board.Grid, board.Territories, h, t, offsetX, offsetY)
		}
	}

	for _, t := range order {
		drawTerritoryDice(screen, t, offsetX, offsetY)
	}
}

func drawTerritoryDice(screen *ebiten.Image, t *game.Territory, offsetX, offsetY float64) {
	cx := offsetX + t.CenterX
	cy := offsetY + t.CenterY

	rightColumnCount := t.NumDice
	if rightColumnCount > game.MaxDice/2 {
		rightColumnCount = game.MaxDice / 2
	}

	leftX := cx - mapDiceSize*mapDiceWidthFactor + mapDiceShiftX
	baseY := cy - mapDiceSize*0.5
	rightX := leftX + mapDiceSize*mapDiceWidthFactor

	for i := game.MaxDice / 2; i < t.NumDice; i++ {
		row := i - game.MaxDice/2
		top := baseY - mapDiceSize*0.3 - float64(row)*mapDiceSize*mapDiceHeightFactor
		drawMapDie(screen, leftX, top, mapDiceSize, t.Owner)
	}
	for i := 0; i < rightColumnCount; i++ {
		top := baseY - float64(i)*mapDiceSize*mapDiceHeightFactor
		drawMapDie(screen, rightX, top, mapDiceSize, t.Owner)
	}
}
