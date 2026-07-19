package app

import (
	"image"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/samuelyuan/dice-wars/internal/game"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	whiteImage.WritePixels(pix)
}

func drawColoredTriangles(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color, antialias bool) {
	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha,
		AntiAlias:      antialias,
	}
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}

func drawFilledPath(dst *ebiten.Image, path *vector.Path, clr color.Color, antialias bool) {
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	drawColoredTriangles(dst, vs, is, clr, antialias)
}

func sortedTerritoriesByDepth(territories []*game.Territory) []*game.Territory {
	out := make([]*game.Territory, 0, len(territories))
	for _, t := range territories {
		if t.Owner >= 0 && len(t.CellIDs) > 0 && t.CenterX >= 0 {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.CenterY != b.CenterY {
			return a.CenterY < b.CenterY
		}
		if a.CenterX != b.CenterX {
			return a.CenterX < b.CenterX
		}
		return a.ID < b.ID
	})
	return out
}

func territoryHighlighted(board *game.Board, t *game.Territory) bool {
	if t.Selected {
		return true
	}
	if board.Phase != game.PhaseDiceRoll && board.Phase != game.PhaseAIAttack {
		return false
	}
	return t.ID == board.SelectedTerr || t.ID == board.OtherTerr
}

func drawOrderTerritories(territories []*game.Territory, board *game.Board) []*game.Territory {
	all := sortedTerritoriesByDepth(territories)
	out := make([]*game.Territory, 0, len(all))
	var highlighted []*game.Territory
	for _, t := range all {
		if territoryHighlighted(board, t) {
			highlighted = append(highlighted, t)
			continue
		}
		out = append(out, t)
	}
	return append(out, highlighted...)
}

func drawHexFill(screen *ebiten.Image, h *game.Hex, offsetX, offsetY float64, fill color.Color) {
	cx := h.Center.X + offsetX
	cy := h.Center.Y + offsetY
	verts := game.HexVertices(cx, cy, h.Radius)

	var path vector.Path
	path.MoveTo(float32(verts[0].X), float32(verts[0].Y))
	for i := 1; i < 6; i++ {
		path.LineTo(float32(verts[i].X), float32(verts[i].Y))
	}
	path.Close()
	drawFilledPath(screen, &path, fill, false)
}

func drawHexBorders(screen *ebiten.Image, board *game.Board, grid *game.HexGrid, territories []*game.Territory, h *game.Hex, t *game.Territory, offsetX, offsetY float64) {
	cx := h.Center.X + offsetX
	cy := h.Center.Y + offsetY
	verts := game.HexVertices(cx, cy, h.Radius)

	highlighted := territoryHighlighted(board, t)
	border := color.RGBA{0, 0, 0, 255}
	if highlighted {
		border = color.RGBA{255, 0, 0, 255}
	}

	for i := 0; i < 6; i++ {
		nb := grid.Neighbor(h, i)
		if nb != nil && nb.TerrID == t.ID {
			continue
		}
		// Shared edges are drawn by the highlighted territory only.
		if !highlighted && nb != nil && nb.TerrID >= 0 && territoryHighlighted(board, territories[nb.TerrID]) {
			continue
		}
		vector.StrokeLine(
			screen,
			float32(verts[i].X), float32(verts[i].Y),
			float32(verts[i+1].X), float32(verts[i+1].Y),
			2, border, false,
		)
	}
}
