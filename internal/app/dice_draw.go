package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/samuelyuan/dice-wars/internal/assets"
)

func drawPlayerDie(screen *ebiten.Image, x, y, size float64, value int, playerIdx int) {
	img := assets.DiceImage(playerIdx, value)
	if img == nil || !drawScaledImage(screen, img, x, y, size) {
		drawPlayerDieFallback(screen, x, y, size, value, playerIdx)
	}
}

func drawMapDie(screen *ebiten.Image, x, y, size float64, playerIdx int) {
	drawPlayerDie(screen, x, y, size, assets.MapDieValue(playerIdx), playerIdx)
}

func drawPlayerIcon(screen *ebiten.Image, centerX, centerY, size float64, playerIdx int) {
	drawMapDie(screen, centerX-size/2, centerY-size/2, size, playerIdx)
}

func drawScaledImage(screen *ebiten.Image, img *ebiten.Image, x, y, size float64) bool {
	width := float64(img.Bounds().Dx())
	if width == 0 {
		return false
	}
	scale := size / width
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	screen.DrawImage(img, op)
	return true
}

func drawPlayerDieFallback(screen *ebiten.Image, x, y, size float64, value int, playerIdx int) {
	fill := lightenPlayerColor(playerIdx, 30)
	stroke := playerStrokeColor(playerIdx)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(size), float32(size), fill, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(size), float32(size), 2, stroke, false)
	drawPips(screen, x, y, size, value, color.RGBA{255, 255, 255, 255})
}

func drawPips(screen *ebiten.Image, x, y, size float64, value int, pipColor color.Color) {
	if value < 1 || value > 6 {
		return
	}

	pipRadius := float32(size * 0.11)
	centerX := float32(x + size/2)
	centerY := float32(y + size/2)
	spacing := float32(size * 0.25)
	r, g, b, a := pipColor.RGBA()
	c := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}

	for _, pip := range diePipOffsets[value-1] {
		vector.DrawFilledCircle(screen, centerX+pip.dx*spacing, centerY+pip.dy*spacing, pipRadius, c, false)
	}
}

type pipOffset struct {
	dx, dy float32
}

var diePipOffsets = [6][]pipOffset{
	{{0, 0}},
	{{-1, -1}, {1, 1}},
	{{-1, -1}, {0, 0}, {1, 1}},
	{{-1, -1}, {1, -1}, {-1, 1}, {1, 1}},
	{{-1, -1}, {1, -1}, {0, 0}, {-1, 1}, {1, 1}},
	{{-1, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {1, 1}},
}
