package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/samuelyuan/dice-wars/internal/game"
)

type Menu struct {
	NumPlayers  int
	ShouldStart bool
}

func NewMenu() *Menu {
	return &Menu{NumPlayers: game.MaxPlayers}
}

// HumanList returns player 1 as human and all others as CPU.
func (m *Menu) HumanList() []bool {
	list := make([]bool, m.NumPlayers)
	list[0] = true
	return list
}

func (m *Menu) Update(lc *LayoutContext) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	mx, my := ebiten.CursorPosition()
	m.handleClick(mx, my, lc)
}

func (m *Menu) handleClick(mx, my int, lc *LayoutContext) {
	if lc.MenuBtnStart().Contains(mx, my) {
		m.ShouldStart = true
		return
	}
	if lc.MenuBtnRemovePlayer().Contains(mx, my) {
		m.adjustPlayerCount(-1)
		return
	}
	if lc.MenuBtnAddPlayer().Contains(mx, my) {
		m.adjustPlayerCount(1)
	}
}

func (m *Menu) adjustPlayerCount(delta int) {
	next := m.NumPlayers + delta
	if next < game.MinPlayers || next > game.MaxPlayers {
		return
	}
	m.NumPlayers = next
}

func (m *Menu) Draw(screen *ebiten.Image, lc *LayoutContext) {
	const titleScale = 6.0
	drawScaledText(screen, "DICE WARS", lc.Width/2, lc.MenuTitleY(), titleScale, colorText)

	summary := m.playersSummary()
	drawText(screen, summary, textCenterX(0, lc.Width, summary), lc.MenuSummaryY(), colorText)

	for i := 0; i < m.NumPlayers; i++ {
		drawMenuPlayerSlot(screen, i, lc)
	}

	lc.MenuBtnRemovePlayer().Draw(screen, false)
	lc.MenuBtnAddPlayer().Draw(screen, false)
	lc.MenuBtnStart().Draw(screen, false)
}

func (m *Menu) playersSummary() string {
	cpuCount := strconv.Itoa(m.NumPlayers - 1)
	return "Players: " + strconv.Itoa(m.NumPlayers) + " (you + " + cpuCount + " CPU)"
}

const (
	menuSlotCircleOffsetX = 20
	menuSlotCircleOffsetY = 40
	menuSlotCircleRadius  = 18
	menuSlotTextZoneW     = 40
	menuSlotLabelOffsetX  = 56
	menuSlotTextY         = 46
)

func drawMenuPlayerSlot(screen *ebiten.Image, playerIdx int, lc *LayoutContext) {
	col := playerIdx % 4
	row := playerIdx / 4
	x := float32(lc.MenuPlayerGridX() + col*lc.MenuPlayerColW())
	y := float32(lc.MenuPlayerGridY() + row*lc.MenuPlayerRowH())

	c := game.PlayerColors[playerIdx]
	circleX := x + menuSlotCircleOffsetX
	circleY := y + menuSlotCircleOffsetY
	vector.DrawFilledCircle(screen, circleX, circleY, menuSlotCircleRadius, color.RGBA{c.R, c.G, c.B, 255}, false)

	textY := int(y) + menuSlotTextY
	num := strconv.Itoa(playerIdx + 1)
	drawText(screen, num, textCenterX(int(x), menuSlotTextZoneW, num), textY, textOnPlayerColor(c.R, c.G, c.B))

	label := "CPU"
	if playerIdx == 0 {
		label = "You"
	}
	drawText(screen, label, int(x)+menuSlotLabelOffsetX, textY, colorText)
}
