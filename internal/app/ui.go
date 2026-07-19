package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Rect struct {
	X, Y, W, H int
}

func (r Rect) Contains(mx, my int) bool {
	return mx >= r.X && mx < r.X+r.W && my >= r.Y && my < r.Y+r.H
}

type Button struct {
	X, Y, W, H int
	Label      string
}

func (b Button) Rect() Rect {
	return Rect{X: b.X, Y: b.Y, W: b.W, H: b.H}
}

func (b Button) Contains(mx, my int) bool {
	return b.Rect().Contains(mx, my)
}

func (b Button) Draw(screen *ebiten.Image, hover bool) {
	bg, fg := buttonColors(hover)
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), bg, false)
	vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 3, color.Black, false)
	drawText(screen, b.Label, textCenterX(b.X, b.W, b.Label), textCenterY(b.Y, b.H), fg)
}

func buttonColors(hover bool) (bg, fg color.RGBA) {
	if hover {
		return color.RGBA{30, 30, 30, 255}, colorTextLight
	}
	return color.RGBA{255, 255, 255, 255}, colorText
}
