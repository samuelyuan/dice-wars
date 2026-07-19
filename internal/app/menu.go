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

func (m *Menu) Update() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	mx, my := ebiten.CursorPosition()
	m.handleClick(mx, my)
}

func (m *Menu) handleClick(mx, my int) {
	if MenuBtnStart.Contains(mx, my) {
		m.ShouldStart = true
		return
	}
	if MenuBtnRemovePlayer.Contains(mx, my) {
		m.adjustPlayerCount(-1)
		return
	}
	if MenuBtnAddPlayer.Contains(mx, my) {
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

func (m *Menu) Draw(screen *ebiten.Image) {
	const titleScale = 6.0
	drawScaledText(screen, "DICE WARS", ScreenWidth/2, 110, titleScale, colorText)

	summary := m.playersSummary()
	drawText(screen, summary, textCenterX(0, ScreenWidth, summary), 190, colorText)

	for i := 0; i < m.NumPlayers; i++ {
		drawMenuPlayerSlot(screen, i)
	}

	MenuBtnRemovePlayer.Draw(screen, false)
	MenuBtnAddPlayer.Draw(screen, false)
	MenuBtnStart.Draw(screen, false)
}

func (m *Menu) playersSummary() string {
	cpuCount := strconv.Itoa(m.NumPlayers - 1)
	return "Players: " + strconv.Itoa(m.NumPlayers) + " (you + " + cpuCount + " CPU)"
}

func drawMenuPlayerSlot(screen *ebiten.Image, playerIdx int) {
	col := playerIdx % 4
	row := playerIdx / 4
	x := float32(MenuPlayerGridX + col*MenuPlayerColW)
	y := float32(MenuPlayerGridY + row*MenuPlayerRowH)

	c := game.PlayerColors[playerIdx]
	vector.DrawFilledCircle(screen, x+20, y+40, 18, color.RGBA{c.R, c.G, c.B, 255}, false)

	num := strconv.Itoa(playerIdx + 1)
	drawText(screen, num, textCenterX(int(x), 40, num), int(y)+46, textOnPlayerColor(c.R, c.G, c.B))

	label := "CPU"
	if playerIdx == 0 {
		label = "You"
	}
	drawText(screen, label, int(x)+56, int(y)+46, colorText)
}
