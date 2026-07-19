package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed pixmaps/*.png
var pixmapsFS embed.FS

var (
	diceOnce  sync.Once
	diceCache [gameMaxPlayers][6]*ebiten.Image
	diceErr   error
)

const gameMaxPlayers = 8

func initDice() {
	for player := 0; player < gameMaxPlayers; player++ {
		for value := 1; value <= 6; value++ {
			img, err := loadDieImage(player, value)
			if err != nil {
				diceErr = err
				return
			}
			diceCache[player][value-1] = img
		}
	}
}

func loadDieImage(player, value int) (*ebiten.Image, error) {
	path := fmt.Sprintf("pixmaps/Player%d_Dice%d.png", player, value)
	data, err := pixmapsFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	processed := stripBlackBackground(decoded)
	trimmed := trimTransparentBounds(processed)
	return ebiten.NewImageFromImage(trimmed), nil
}

// DiceImage returns the die face for a player (0-7) and value (1-6).
func DiceImage(playerIdx, value int) *ebiten.Image {
	diceOnce.Do(initDice)
	if diceErr != nil || playerIdx < 0 || playerIdx >= gameMaxPlayers || value < 1 || value > 6 {
		return nil
	}
	return diceCache[playerIdx][value-1]
}

// MapDieValue picks the stacked-die face shown on territories.
func MapDieValue(playerIdx int) int {
	return playerIdx%6 + 1
}

func stripBlackBackground(src image.Image) image.Image {
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBAModel.Convert(src.At(x, y)).(color.RGBA)
			if isNearBlack(c) {
				c.A = 0
			}
			out.Set(x, y, c)
		}
	}
	return out
}

func isNearBlack(c color.RGBA) bool {
	return c.R < 24 && c.G < 24 && c.B < 24
}

func trimTransparentBounds(src image.Image) image.Image {
	bounds := src.Bounds()
	crop, ok := opaqueBounds(src, bounds)
	if !ok {
		return src
	}
	return copyRegion(src, crop)
}

func opaqueBounds(src image.Image, bounds image.Rectangle) (image.Rectangle, bool) {
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y
	found := false

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if _, _, _, alpha := src.At(x, y).RGBA(); alpha == 0 {
				continue
			}
			found = true
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}
	}
	if !found {
		return bounds, false
	}
	return image.Rect(minX, minY, maxX+1, maxY+1), true
}

func copyRegion(src image.Image, crop image.Rectangle) image.Image {
	out := image.NewRGBA(image.Rect(0, 0, crop.Dx(), crop.Dy()))
	for y := crop.Min.Y; y < crop.Max.Y; y++ {
		for x := crop.Min.X; x < crop.Max.X; x++ {
			out.Set(x-crop.Min.X, y-crop.Min.Y, src.At(x, y))
		}
	}
	return out
}
