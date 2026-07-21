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

// BtnFastForward toggles a large speed multiplier for blowing through
// AI-heavy stretches.
func (lc *LayoutContext) BtnFastForward() Button {
	return Button{X: 140, Y: int(lc.PlayerBarY()) - 60, W: 100, H: 30, Label: "Fast Fwd"}
}

func (lc *LayoutContext) BtnEndTurn() Button {
	return Button{X: lc.Width - 140, Y: int(lc.PlayerBarY()) - 50, W: 120, H: 40, Label: "End Turn"}
}

func (lc *LayoutContext) BtnAuto() Button {
	return Button{X: lc.Width - 270, Y: int(lc.PlayerBarY()) - 50, W: 100, H: 40, Label: "Auto"}
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

// Replay control bar: sits below the player bar, replacing the live-game button row.
const (
	replayBarHeight   = 50
	replayPlayBtnSize = 40
	replaySpeedBtnW   = 46
	replayExitBtnW    = 70
	replaySideBtnH    = 30
	replayBtnGap      = 10
	replaySeekBarMinW = 100
)

func (lc *LayoutContext) ReplayBarY() int {
	return int(lc.PlayerBarY()) + PlayerBarHeight + 8
}

func (lc *LayoutContext) ReplayPlayButton() Button {
	y := lc.ReplayBarY() + (replayBarHeight-replayPlayBtnSize)/2
	return Button{X: 20, Y: y, W: replayPlayBtnSize, H: replayPlayBtnSize}
}

// replaySpeedOptions are the playback speeds shown as chip buttons.
var replaySpeedOptions = []float64{0.5, 1, 2, 3}

func (lc *LayoutContext) ReplaySpeedButtons() []Button {
	totalW := len(replaySpeedOptions)*replaySpeedBtnW + (len(replaySpeedOptions)-1)*replayBtnGap
	startX := lc.Width - replayBtnGap - replayExitBtnW - replayBtnGap - totalW
	y := lc.ReplayBarY() + (replayBarHeight-replaySideBtnH)/2

	btns := make([]Button, len(replaySpeedOptions))
	for i, speed := range replaySpeedOptions {
		label := formatSpeed(speed) + "x"
		btns[i] = Button{X: startX + i*(replaySpeedBtnW+replayBtnGap), Y: y, W: replaySpeedBtnW, H: replaySideBtnH, Label: label}
	}
	return btns
}

func (lc *LayoutContext) ReplayExitButton() Button {
	y := lc.ReplayBarY() + (replayBarHeight-replaySideBtnH)/2
	return Button{X: lc.Width - replayBtnGap - replayExitBtnW, Y: y, W: replayExitBtnW, H: replaySideBtnH, Label: "Exit"}
}

// ReplaySeekBar spans between the play button and the right-side controls.
func (lc *LayoutContext) ReplaySeekBar() Rect {
	play := lc.ReplayPlayButton()
	x := play.X + play.W + 20
	speedBtns := lc.ReplaySpeedButtons()
	rightEdge := speedBtns[0].X - 20
	w := rightEdge - x
	if w < replaySeekBarMinW {
		w = replaySeekBarMinW
	}
	y := lc.ReplayBarY() + replayBarHeight/2 - 4
	return Rect{X: x, Y: y, W: w, H: 8}
}

// Three-button row shared by the elimination overlay and the victory screen.
const (
	eliminationBtnW   = 180
	eliminationBtnH   = 50
	eliminationBtnGap = 20
)

func (lc *LayoutContext) eliminationButtonsY() int {
	return lc.Height/2 + 30
}

func (lc *LayoutContext) threeButtonRowStartX() int {
	total := 3*eliminationBtnW + 2*eliminationBtnGap
	return lc.Width/2 - total/2
}

// threeButtonRowButton returns the button at the given slot (0-2) in a
// three-button row, shared by the elimination overlay and victory screen.
func (lc *LayoutContext) threeButtonRowButton(slot int, y int, label string) Button {
	x := lc.threeButtonRowStartX() + slot*(eliminationBtnW+eliminationBtnGap)
	return Button{X: x, Y: y, W: eliminationBtnW, H: eliminationBtnH, Label: label}
}

func (lc *LayoutContext) EliminationReplayButton() Button {
	return lc.threeButtonRowButton(0, lc.eliminationButtonsY(), "Replay")
}

func (lc *LayoutContext) EliminationRestartButton() Button {
	return lc.threeButtonRowButton(1, lc.eliminationButtonsY(), "Restart")
}

func (lc *LayoutContext) EliminationNewGameButton() Button {
	return lc.threeButtonRowButton(2, lc.eliminationButtonsY(), "New Game")
}

func (lc *LayoutContext) victoryButtonsY() int {
	return lc.Height/2 + 60
}

func (lc *LayoutContext) VictoryReplayButton() Button {
	return lc.threeButtonRowButton(0, lc.victoryButtonsY(), "Replay")
}

func (lc *LayoutContext) VictoryRestartButton() Button {
	return lc.threeButtonRowButton(1, lc.victoryButtonsY(), "Restart")
}

func (lc *LayoutContext) VictoryMenuButton() Button {
	return lc.threeButtonRowButton(2, lc.victoryButtonsY(), "Main Menu")
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
