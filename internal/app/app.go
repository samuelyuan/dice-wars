package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/samuelyuan/dice-wars/internal/game"
)

const frameDuration = 1.0 / 60.0

type App struct {
	screen   Screen
	menu     *Menu
	board    *game.Board
	hoverBtn string
	wantMenu bool
	layout   *LayoutContext
}

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenGame
	ScreenVictory
)

func NewApp() *App {
	return &App{
		screen: ScreenMenu,
		menu:   NewMenu(),
		layout: &LayoutContext{Width: DefaultScreenWidth, Height: DefaultScreenHeight},
	}
}

func (a *App) Update() error {
	// Update layout context with current window size
	w, h := ebiten.WindowSize()
	a.layout.Width = w
	a.layout.Height = h

	// Toggle fullscreen with F11
	if ebiten.IsKeyPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	switch a.screen {
	case ScreenMenu:
		a.updateMenu()
	case ScreenGame:
		a.updateGame()
	case ScreenVictory:
		a.updateVictory()
	}
	return nil
}

func (a *App) updateMenu() {
	a.menu.Update(a.layout)
	if !a.menu.ShouldStart {
		return
	}
	a.menu.ShouldStart = false
	a.board = game.NewBoard(a.menu.NumPlayers, a.menu.HumanList())
	a.screen = ScreenGame
}

func (a *App) updateGame() {
	if a.board == nil {
		return
	}
	a.handleGameInput()
	if a.wantMenu {
		a.wantMenu = false
		a.board = nil
		a.screen = ScreenMenu
		return
	}
	a.board.Update(frameDuration)
	if a.board.GameOver {
		a.screen = ScreenVictory
	}
}

func (a *App) updateVictory() {
	if !ebiten.IsKeyPressed(ebiten.KeyEnter) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return
	}
	a.board = nil
	a.screen = ScreenMenu
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	switch a.screen {
	case ScreenMenu:
		a.menu.Draw(screen, a.layout)
	case ScreenGame:
		if a.board != nil {
			drawGame(screen, a.board, a.hoverBtn, a.layout)
		}
	case ScreenVictory:
		if a.board != nil {
			drawVictory(screen, a.board.VictoryPlayer, a.board.VictoryHuman, a.layout)
		}
	}
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func drawVictory(screen *ebiten.Image, player int, human bool, lc *LayoutContext) {
	msg := "Player " + strconv.Itoa(player+1) + " wins!"
	if human {
		msg = "You win!"
	}
	centerY := lc.Height / 2
	drawText(screen, "VICTORY!", textCenterX(0, lc.Width, "VICTORY!"), centerY-20, colorText)
	drawText(screen, msg, textCenterX(0, lc.Width, msg), centerY+20, colorText)
	hint := "Click or press Enter for menu"
	drawText(screen, hint, textCenterX(0, lc.Width, hint), centerY+50, colorText)
}
