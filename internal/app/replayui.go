package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorSeekTrack  = color.RGBA{200, 200, 200, 255}
	colorSeekFill   = color.RGBA{60, 60, 60, 255}
	colorSeekHandle = color.RGBA{20, 20, 20, 255}
)

// formatSpeed drops a trailing ".0" (1 instead of 1.0; 0.5 stays as-is).
func formatSpeed(speed float64) string {
	s := strconv.FormatFloat(speed, 'f', -1, 64)
	return s
}

// handleReplayInput processes clicks on the video-player control bar.
func (a *App) handleReplayInput() {
	if a.replayPlayer == nil {
		return
	}
	lc := a.layout
	mx, my := ebiten.CursorPosition()

	seekBar := lc.ReplaySeekBar()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && seekBarHitArea(seekBar).Contains(mx, my) {
		a.replayDragging = true
	}
	if a.replayDragging {
		a.replaySeekPreview = seekFraction(mx, seekBar)
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			a.replayDragging = false
			a.replayPlayer.SeekToProgress(a.replaySeekPreview)
		}
		return
	}

	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	if lc.ReplayPlayButton().Contains(mx, my) {
		if a.replayIsFinished() {
			a.replayPlayer.SeekToProgress(0)
			a.replayPaused = false
			return
		}
		a.replayPaused = !a.replayPaused
		return
	}
	if lc.ReplayExitButton().Contains(mx, my) {
		a.replayPlayer = nil
		a.screen = ScreenMenu
		return
	}
	for i, btn := range lc.ReplaySpeedButtons() {
		if btn.Contains(mx, my) {
			a.replaySpeed = replaySpeedOptions[i]
			return
		}
	}
}

// seekBarHitArea pads the thin seek track so it's easier to click/grab.
func seekBarHitArea(bar Rect) Rect {
	const pad = 10
	return Rect{X: bar.X, Y: bar.Y - pad, W: bar.W, H: bar.H + pad*2}
}

func seekFraction(mx int, bar Rect) float64 {
	if bar.W == 0 {
		return 0
	}
	f := float64(mx-bar.X) / float64(bar.W)
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	return f
}

func (a *App) drawReplayControls(screen *ebiten.Image, lc *LayoutContext) {
	finished := a.replayIsFinished()
	mx, my := ebiten.CursorPosition()

	playBtn := lc.ReplayPlayButton()
	drawReplayPlayButton(screen, playBtn, a.replayPaused, finished, playBtn.Contains(mx, my))
	drawReplaySeekBar(screen, lc.ReplaySeekBar(), a.replaySeekProgress())
	drawReplaySpeedButtons(screen, lc.ReplaySpeedButtons(), a.replaySpeed, mx, my)
	exitBtn := lc.ReplayExitButton()
	exitBtn.Draw(screen, exitBtn.Contains(mx, my))

	statusText := "Playing"
	switch {
	case finished:
		statusText = "Finished"
	case a.replayPaused:
		statusText = "Paused"
	}
	label := statusText + "  Move " + strconv.Itoa(a.replayPlayer.ActionsPlayed()) + "/" + strconv.Itoa(len(a.lastReplay.Actions))
	drawText(screen, label, lc.ReplayPlayButton().X, lc.ReplayBarY()-14, colorText)
}

// replayIsFinished mirrors updateReplay's auto-stop so the label/icon match.
func (a *App) replayIsFinished() bool {
	return (a.lastReplay.IsPartialGame() && a.replayPlayer.Board.HumanEliminated()) ||
		a.replayPlayer.Finished()
}

// replaySeekProgress returns the drag preview while scrubbing, else the actual position.
func (a *App) replaySeekProgress() float64 {
	if a.replayDragging {
		return a.replaySeekPreview
	}
	return a.replayPlayer.Progress()
}

// drawReplayPlayButton shows a restart icon when finished, else play/pause.
func drawReplayPlayButton(screen *ebiten.Image, btn Button, paused, finished, hover bool) {
	bg, fg := buttonColors(hover)
	vector.DrawFilledRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), bg, false)
	vector.StrokeRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), 3, color.Black, false)

	cx := float64(btn.X + btn.W/2)
	cy := float64(btn.Y + btn.H/2)
	const iconSize = 12.0

	switch {
	case finished:
		drawRestartIcon(screen, cx, cy, iconSize, fg)
	case paused:
		drawPlayTriangle(screen, cx, cy, iconSize, fg)
	default:
		drawPauseBars(screen, cx, cy, iconSize, fg)
	}
}

func drawPlayTriangle(screen *ebiten.Image, cx, cy, size float64, clr color.Color) {
	var path vector.Path
	path.MoveTo(float32(cx-size*0.4), float32(cy-size))
	path.LineTo(float32(cx-size*0.4), float32(cy+size))
	path.LineTo(float32(cx+size*0.8), float32(cy))
	path.Close()
	drawFilledPath(screen, &path, clr, true)
}

func drawPauseBars(screen *ebiten.Image, cx, cy, size float64, clr color.Color) {
	barW := float32(size * 0.5)
	barH := float32(size * 2)
	gap := float32(size * 0.5)
	vector.DrawFilledRect(screen, float32(cx)-gap/2-barW, float32(cy)-barH/2, barW, barH, clr, false)
	vector.DrawFilledRect(screen, float32(cx)+gap/2, float32(cy)-barH/2, barW, barH, clr, false)
}

// drawRestartIcon draws the standard "skip to start" bar+triangle glyph.
func drawRestartIcon(screen *ebiten.Image, cx, cy, size float64, clr color.Color) {
	barW := float32(size * 0.3)
	barH := float32(size * 2)
	barX := float32(cx - size*0.9)
	vector.DrawFilledRect(screen, barX, float32(cy)-barH/2, barW, barH, clr, false)

	var tri vector.Path
	tri.MoveTo(float32(cx+size*0.4), float32(cy-size))
	tri.LineTo(float32(cx+size*0.4), float32(cy+size))
	tri.LineTo(float32(cx-size*0.4), float32(cy))
	tri.Close()
	drawFilledPath(screen, &tri, clr, true)
}

func drawReplaySeekBar(screen *ebiten.Image, bar Rect, progress float64) {
	vector.DrawFilledRect(screen, float32(bar.X), float32(bar.Y), float32(bar.W), float32(bar.H), colorSeekTrack, false)

	fillW := float32(bar.W) * float32(progress)
	if fillW > 0 {
		vector.DrawFilledRect(screen, float32(bar.X), float32(bar.Y), fillW, float32(bar.H), colorSeekFill, false)
	}

	handleX := float32(bar.X) + fillW
	handleY := float32(bar.Y) + float32(bar.H)/2
	vector.DrawFilledCircle(screen, handleX, handleY, 8, colorSeekHandle, true)
}

func drawReplaySpeedButtons(screen *ebiten.Image, btns []Button, activeSpeed float64, mx, my int) {
	for i, btn := range btns {
		active := replaySpeedOptions[i] == activeSpeed
		bg, fg := buttonColors(active || btn.Contains(mx, my))
		vector.DrawFilledRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), bg, false)
		vector.StrokeRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), 2, color.Black, false)
		drawText(screen, btn.Label, textCenterX(btn.X, btn.W, btn.Label), textCenterY(btn.Y, btn.H), fg)
	}
}
