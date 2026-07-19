package game

import "testing"

func TestPickHexAtCenter(t *testing.T) {
	g := NewHexGrid()
	for ax, h := range g.Hexes {
		picked := g.PickHex(h.Center.X, h.Center.Y)
		if picked == nil {
			t.Fatalf("nil pick at center of %v (%.2f, %.2f)", ax, h.Center.X, h.Center.Y)
		}
		if picked.Axial != ax {
			t.Fatalf("pick at (%.2f,%.2f) want %v got %v", h.Center.X, h.Center.Y, ax, picked.Axial)
		}
	}
}
