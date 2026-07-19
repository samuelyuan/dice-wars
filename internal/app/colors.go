package app

import (
	"image/color"

	"github.com/samuelyuan/dice-wars/internal/game"
)

func playerFill(playerIdx int, selected bool) color.RGBA {
	if selected {
		return color.RGBA{0, 0, 0, 255}
	}
	return scalePlayerColor(playerIdx, 1.4)
}

func scalePlayerColor(playerIdx int, scale float64) color.RGBA {
	c := game.PlayerColors[playerIdx]
	return color.RGBA{
		R: clampByte(int(float64(c.R) * scale)),
		G: clampByte(int(float64(c.G) * scale)),
		B: clampByte(int(float64(c.B) * scale)),
		A: 255,
	}
}

func lightenPlayerColor(playerIdx int, delta int) color.RGBA {
	c := game.PlayerColors[playerIdx]
	return color.RGBA{
		R: clampByte(int(c.R) + delta),
		G: clampByte(int(c.G) + delta),
		B: clampByte(int(c.B) + delta),
		A: 255,
	}
}

func playerStrokeColor(playerIdx int) color.RGBA {
	c := game.PlayerColors[playerIdx]
	return color.RGBA{c.R, c.G, c.B, 255}
}

func clampByte(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
