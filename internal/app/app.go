package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/samuelyuan/dice-wars/internal/game"
)

const frameDuration = 1.0 / 60.0

type App struct {
	screen                 Screen
	menu                   *Menu
	board                  *game.Board
	hoverBtn               string
	wantMenu               bool
	layout                 *LayoutContext
	fastForward            bool    // speeds up AI turns while true
	humanEliminated        bool    // true once we've detected the human losing all territory this game
	showEliminationOverlay bool    // whether the "you've been eliminated" overlay is up
	lastReplay             *game.Replay
	replayPlayer           *game.ReplayPlayer
	replaySpeed            float64 // Playback speed multiplier
	replayPaused           bool
	replayDragging         bool    // scrubbing the seek bar
	replaySeekPreview      float64 // live drag position while scrubbing, [0,1]
}

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenGame
	ScreenVictory
	ScreenReplay
)

func NewApp() *App {
	return &App{
		screen:      ScreenMenu,
		menu:        NewMenu(),
		layout:      &LayoutContext{Width: DefaultScreenWidth, Height: DefaultScreenHeight},
		replaySpeed: 1.0,
	}
}

// fastForwardSpeed is the multiplier applied to AI turns while Fast Forward
// is toggled on.
const fastForwardSpeed = 50.0

func (a *App) effectiveGameSpeed() float64 {
	if a.fastForward {
		return fastForwardSpeed
	}
	return 1.0
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
	case ScreenReplay:
		a.updateReplay()
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
	a.humanEliminated = false
	a.showEliminationOverlay = false
	a.screen = ScreenGame
}

func (a *App) updateGame() {
	if a.board == nil {
		return
	}
	if a.showEliminationOverlay {
		a.handleEliminationOverlayInput()
		return
	}
	a.handleGameInput()
	if a.wantMenu {
		a.wantMenu = false
		a.board = nil
		a.humanEliminated = false
		a.screen = ScreenMenu
		return
	}
	a.board.Update(frameDuration * a.effectiveGameSpeed())
	if a.board.GameOver {
		a.lastReplay = a.board.ExportReplay()
		a.screen = ScreenVictory
		return
	}
	if !a.humanEliminated && a.board.HumanEliminated() {
		a.humanEliminated = true
		a.showEliminationOverlay = true
	}
}

func (a *App) updateVictory() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		lc := a.layout
		mx, my := ebiten.CursorPosition()
		switch {
		case lc.VictoryReplayButton().Contains(mx, my):
			a.startReplay()
			return
		case lc.VictoryRestartButton().Contains(mx, my):
			a.board = game.NewBoard(a.menu.NumPlayers, a.menu.HumanList())
			a.screen = ScreenGame
			return
		case lc.VictoryMenuButton().Contains(mx, my):
			a.board = nil
			a.screen = ScreenMenu
			return
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		a.board = nil
		a.screen = ScreenMenu
	}
}

func (a *App) startReplay() {
	if a.lastReplay == nil {
		return
	}
	a.replayPlayer = game.NewReplayPlayer(a.lastReplay)
	a.replayPaused = false
	a.screen = ScreenReplay
}

func (a *App) updateReplay() {
	if a.replayPlayer == nil {
		a.screen = ScreenMenu
		return
	}

	a.handleReplayInput()
	if a.replayPlayer == nil {
		// handleReplayInput may have exited replay mode (e.g. Exit button).
		return
	}

	// Keyboard shortcuts alongside the on-screen controls
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		a.replayPaused = !a.replayPaused
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		a.stepReplaySpeed(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		a.stepReplaySpeed(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		a.replayPlayer = nil
		a.screen = ScreenMenu
		return
	}

	// Skip while dragging so the preview doesn't fight live playback.
	if !a.replayPaused && !a.replayDragging {
		a.replayPlayer.Board.Update(frameDuration * a.replaySpeed)
		if a.replayIsFinished() {
			a.replayPaused = true
		}
	}
}

func (a *App) stepReplaySpeed(delta int) {
	idx := 0
	for i, s := range replaySpeedOptions {
		if s == a.replaySpeed {
			idx = i
			break
		}
	}
	idx += delta
	if idx < 0 {
		idx = 0
	}
	if idx >= len(replaySpeedOptions) {
		idx = len(replaySpeedOptions) - 1
	}
	a.replaySpeed = replaySpeedOptions[idx]
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	switch a.screen {
	case ScreenMenu:
		a.menu.Draw(screen, a.layout)
	case ScreenGame:
		if a.board != nil {
			drawGame(screen, a.board, a.hoverBtn, a.layout, true, a.fastForward)
			if a.showEliminationOverlay {
				a.drawEliminationOverlay(screen, a.layout)
			}
		}
	case ScreenVictory:
		a.drawVictory(screen, a.layout)
	case ScreenReplay:
		if a.replayPlayer != nil {
			drawGame(screen, a.replayPlayer.Board, "", a.layout, false, false)
			a.drawReplayControls(screen, a.layout)
		}
	}
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (a *App) drawVictory(screen *ebiten.Image, lc *LayoutContext) {
	if a.board == nil {
		return
	}
	header := "VICTORY!"
	msg := "Player " + strconv.Itoa(a.board.VictoryPlayer+1) + " wins!"
	switch {
	case a.board.VictoryHuman:
		msg = "You win!"
	case a.board.HumanIndex() >= 0:
		// GameOver can fire on the same attack that eliminates the human,
		// skipping the elimination overlay — still show the loss here.
		header = "DEFEAT"
	}
	centerY := lc.Height / 2
	drawText(screen, header, textCenterX(0, lc.Width, header), centerY-20, colorText)
	drawText(screen, msg, textCenterX(0, lc.Width, msg), centerY+20, colorText)

	mx, my := ebiten.CursorPosition()
	lc.VictoryReplayButton().Draw(screen, lc.VictoryReplayButton().Contains(mx, my))
	lc.VictoryRestartButton().Draw(screen, lc.VictoryRestartButton().Contains(mx, my))
	lc.VictoryMenuButton().Draw(screen, lc.VictoryMenuButton().Contains(mx, my))
}
