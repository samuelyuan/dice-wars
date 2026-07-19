package app

import "github.com/samuelyuan/dice-wars/internal/game"

// LayoutContext holds screen dimensions for dynamic layout calculations
type LayoutContext struct {
	Width  int
	Height int
}

// Default layout constants (minimum/recommended size)
const (
	DefaultScreenWidth  = 1280
	DefaultScreenHeight = 820
)

const (
	PlayerBarHeight = 40
	PlayerBarSlotW  = 120
)

const (
	TurnBannerW = 400
	TurnBannerH = 32
)

const (
	minBattlePanelHeight = 32.0
	minBattleDiceSize    = 18.0
	battlePanelPadding   = 8.0
)

// MapOffsetX centers the hex map horizontally within the window.
func (lc *LayoutContext) MapOffsetX() float64 {
	offset := (float64(lc.Width) - game.MapPixelWidth()) / 2
	if offset < 10 {
		offset = 10
	}
	return offset
}

// MapOffsetY anchors the hex map just below the turn banner.
func (lc *LayoutContext) MapOffsetY() float64 {
	return lc.TurnBannerY() + TurnBannerH + 10
}

// Layout calculation methods
func (lc *LayoutContext) PlayerBarY() float64 {
	return float64(lc.Height) - 110
}

func (lc *LayoutContext) StatusTextY() float64 {
	return float64(lc.Height) - 24
}

func (lc *LayoutContext) TurnBannerX() float64 {
	return float64(lc.Width)/2 - TurnBannerW/2
}

func (lc *LayoutContext) TurnBannerY() float64 {
	return 8
}

func (lc *LayoutContext) BtnMenu() Button {
	return Button{X: 10, Y: int(lc.PlayerBarY()) - 60, W: 120, H: 30, Label: "Menu"}
}

func (lc *LayoutContext) BtnEndTurn() Button {
	return Button{X: lc.Width - 140, Y: int(lc.PlayerBarY()) - 50, W: 120, H: 40, Label: "End Turn"}
}

func (lc *LayoutContext) BtnAuto() Button {
	return Button{X: lc.Width - 270, Y: int(lc.PlayerBarY()) - 50, W: 100, H: 40, Label: "Auto"}
}

func (lc *LayoutContext) CheatHit() Rect {
	return Rect{X: 300, Y: 5, W: 50, H: 20}
}

func (lc *LayoutContext) MenuBtnStart() Button {
	return Button{X: lc.Width - 220, Y: lc.Height - 80, W: 200, H: 60, Label: "Start!"}
}

func (lc *LayoutContext) MenuBtnRemovePlayer() Button {
	return Button{X: lc.Width/2 - 70, Y: lc.Height - 110, W: 60, H: 60, Label: "-"}
}

func (lc *LayoutContext) MenuBtnAddPlayer() Button {
	return Button{X: lc.Width/2 + 10, Y: lc.Height - 110, W: 60, H: 60, Label: "+"}
}

const menuPlayerGridCols = 4

func (lc *LayoutContext) MenuTitleY() int {
	return lc.Height / 8
}

func (lc *LayoutContext) MenuSummaryY() int {
	return lc.MenuTitleY() + 70
}

func (lc *LayoutContext) MenuPlayerGridX() int {
	gridWidth := menuPlayerGridCols * lc.MenuPlayerColW()
	return (lc.Width - gridWidth) / 2
}

func (lc *LayoutContext) MenuPlayerGridY() int {
	y := lc.MenuSummaryY() + 50
	const maxRows = 2
	maxBottom := lc.Height - 130 - maxRows*lc.MenuPlayerRowH()
	if y > maxBottom {
		y = maxBottom
	}
	return y
}

func (lc *LayoutContext) MenuPlayerColW() int {
	return 200
}

func (lc *LayoutContext) MenuPlayerRowH() int {
	return 90
}

type battlePanelLayout struct {
	centerY  float64
	height   float64
	diceSize float64
}

func (lc *LayoutContext) mapContentBottom() float64 {
	rowStep := game.HexRadius * 2 * 3 / 4
	lastRowCenterY := rowStep * float64(game.GridHeight-1)
	return lc.MapOffsetY() + lastRowCenterY + game.HexRadius + mapDiceSize*0.6
}

func (lc *LayoutContext) computeBattlePanelLayout() battlePanelLayout {
	top := lc.mapContentBottom() + battlePanelPadding
	bottom := lc.PlayerBarY() - battlePanelPadding
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
