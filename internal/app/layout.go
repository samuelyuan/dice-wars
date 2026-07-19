package app

import "github.com/samuelyuan/dice-wars/internal/game"

const (
	ScreenWidth  = 1280
	ScreenHeight = 820
	MapOffsetX   = 50.0
	MapOffsetY   = 50.0
)

const (
	PlayerBarY      = ScreenHeight - 110
	PlayerBarHeight = 40
	PlayerBarSlotW  = 120
	StatusTextY     = ScreenHeight - 24
)

const (
	TurnBannerW = 400
	TurnBannerH = 32
	TurnBannerX = ScreenWidth/2 - TurnBannerW/2
	TurnBannerY = 8
)

var (
	BtnMenu             = Button{X: 10, Y: ScreenHeight - 50, W: 120, H: 30, Label: "Menu"}
	BtnEndTurn          = Button{X: ScreenWidth - 140, Y: ScreenHeight - 60, W: 120, H: 40, Label: "End Turn"}
	BtnAuto             = Button{X: ScreenWidth - 270, Y: ScreenHeight - 60, W: 100, H: 40, Label: "Auto"}
	CheatHit            = Rect{X: 300, Y: 5, W: 50, H: 20}
	MenuBtnStart        = Button{X: ScreenWidth - 220, Y: ScreenHeight - 80, W: 200, H: 60, Label: "Start!"}
	MenuBtnRemovePlayer = Button{X: ScreenWidth/2 - 70, Y: ScreenHeight - 110, W: 60, H: 60, Label: "-"}
	MenuBtnAddPlayer    = Button{X: ScreenWidth/2 + 10, Y: ScreenHeight - 110, W: 60, H: 60, Label: "+"}
	MenuPlayerGridX     = 100
	MenuPlayerGridY     = 320
	MenuPlayerColW      = 200
	MenuPlayerRowH      = 90
)

const (
	minBattlePanelHeight = 32.0
	minBattleDiceSize    = 18.0
	battlePanelPadding   = 8.0
)

type battlePanelLayout struct {
	centerY  float64
	height   float64
	diceSize float64
}

func mapContentBottom() float64 {
	rowStep := game.HexRadius * 2 * 3 / 4
	lastRowCenterY := rowStep * float64(game.GridHeight-1)
	return MapOffsetY + lastRowCenterY + game.HexRadius + mapDiceSize*0.6
}

func computeBattlePanelLayout() battlePanelLayout {
	top := mapContentBottom() + battlePanelPadding
	bottom := float64(PlayerBarY) - battlePanelPadding
	height := bottom - top
	if height < minBattlePanelHeight {
		height = minBattlePanelHeight
		top = bottom - height
	}

	diceSize := battleDiceSize
	if height < 64 {
		diceSize = height * 0.52
		if diceSize < minBattleDiceSize {
			diceSize = minBattleDiceSize
		}
	}

	return battlePanelLayout{
		centerY:  top + height/2,
		height:   height,
		diceSize: diceSize,
	}
}
