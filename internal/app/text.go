package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

var (
	colorText      = color.RGBA{25, 25, 25, 255}
	colorTextLight = color.RGBA{255, 255, 255, 255}
)

func drawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	if str == "" {
		return
	}
	text.Draw(screen, str, basicfont.Face7x13, x, y, clr)
}

func drawScaledText(screen *ebiten.Image, str string, centerX, centerY int, scale float64, clr color.Color) {
	if str == "" {
		return
	}
	const charW, charH = 7, 13
	w := len(str) * charW
	img := ebiten.NewImage(w, charH)
	text.Draw(img, str, basicfont.Face7x13, 0, charH, clr)

	scaledW := float64(w) * scale
	scaledH := float64(charH) * scale
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(centerX)-scaledW/2, float64(centerY)-scaledH/2)
	screen.DrawImage(img, op)
}

func textWidth(str string) int {
	return len(str) * 7
}

func textCenterX(x, w int, str string) int {
	return x + (w-textWidth(str))/2
}

func textCenterY(y, h int) int {
	return y + (h+13)/2 - 2
}

func textOnPlayerColor(r, g, b uint8) color.Color {
	lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	if lum > 150 {
		return colorText
	}
	return colorTextLight
}
