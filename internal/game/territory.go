package game

import "math/rand/v2"

type Territory struct {
	ID         int
	Owner      int // player index, -1 if none
	NumDice    int
	Selected   bool
	CellIDs    []Axial
	Neighbours []int // territory IDs
	CenterX    float64
	CenterY    float64
}

func (t *Territory) regenerateNeighbours(grid *HexGrid, territories []*Territory) {
	t.Neighbours = t.Neighbours[:0]
	for _, ax := range t.CellIDs {
		h := grid.Hexes[ax]
		forEachHexNeighbor(grid, h, func(nb *Hex) {
			if nb.TerrID < 0 {
				return
			}
			other := territories[nb.TerrID]
			if other.ID != t.ID && other.Owner >= 0 {
				t.Neighbours = appendIntIfMissing(t.Neighbours, other.ID)
			}
		})
	}
}

func (t *Territory) appendCell(grid *HexGrid, territories []*Territory, h *Hex) {
	if h.TerrID >= 0 {
		t.takeCellFrom(grid, territories, h)
	} else {
		t.linkAdjacentTerritories(grid, territories, h)
	}
	t.CellIDs = append(t.CellIDs, h.Axial)
	h.TerrID = t.ID
	t.calculateCenter(grid)
}

func (t *Territory) takeCellFrom(grid *HexGrid, territories []*Territory, h *Hex) {
	old := territories[h.TerrID]
	old.CellIDs = removeAxial(old.CellIDs, h.Axial)
	old.regenerateNeighbours(grid, territories)
}

func (t *Territory) linkAdjacentTerritories(grid *HexGrid, territories []*Territory, h *Hex) {
	forEachHexNeighbor(grid, h, func(nb *Hex) {
		if nb.TerrID < 0 {
			return
		}
		other := territories[nb.TerrID]
		if other.ID == t.ID || other.Owner < 0 {
			return
		}
		t.Neighbours = appendIntIfMissing(t.Neighbours, other.ID)
		other.Neighbours = appendIntIfMissing(other.Neighbours, t.ID)
	})
}

func removeAxial(cells []Axial, target Axial) []Axial {
	for i, ax := range cells {
		if ax == target {
			return append(cells[:i], cells[i+1:]...)
		}
	}
	return cells
}

func forEachHexNeighbor(grid *HexGrid, h *Hex, fn func(*Hex)) {
	for dir := 0; dir < 6; dir++ {
		if nb := grid.Neighbor(h, dir); nb != nil {
			fn(nb)
		}
	}
}

func (t *Territory) calculateCenter(grid *HexGrid) {
	if len(t.CellIDs) == 0 {
		t.CenterX, t.CenterY = -1, -1
		return
	}
	var sumX, sumY float64
	for _, ax := range t.CellIDs {
		h := grid.Hexes[ax]
		sumX += h.Center.X
		sumY += h.Center.Y
	}
	n := float64(len(t.CellIDs))
	t.CenterX = sumX / n
	t.CenterY = sumY / n
}

func (t *Territory) findEmptyAdjacent(rng *rand.Rand, grid *HexGrid) *Hex {
	if len(t.CellIDs) == 0 {
		return nil
	}
	var found *Hex
	forEachShuffled(len(t.CellIDs), rng, func(cellIdx int) bool {
		hex := grid.Hexes[t.CellIDs[cellIdx]]
		found = findFirstEmptyNeighbor(rng, grid, hex)
		return found != nil
	})
	return found
}

func findFirstEmptyNeighbor(rng *rand.Rand, grid *HexGrid, hex *Hex) *Hex {
	var found *Hex
	forEachShuffled(6, rng, func(dir int) bool {
		nb := grid.Neighbor(hex, dir)
		if nb != nil && nb.TerrID < 0 {
			found = nb
			return true
		}
		return false
	})
	return found
}

func (t *Territory) grow(rng *rand.Rand, grid *HexGrid, territories []*Territory, targetCells int) int {
	if len(t.CellIDs) == 0 {
		return 0
	}

	cellCount := len(t.CellIDs)
	for grown := 0; grown < targetCells; grown++ {
		if !t.growOneCell(rng, grid, territories, &cellCount) {
			return grown
		}
	}
	return targetCells
}

func (t *Territory) growOneCell(rng *rand.Rand, grid *HexGrid, territories []*Territory, cellCount *int) bool {
	added := false
	forEachShuffled(*cellCount, rng, func(cellIdx int) bool {
		hex := grid.Hexes[t.CellIDs[cellIdx]]
		empty := findFirstEmptyNeighbor(rng, grid, hex)
		if empty == nil {
			return false
		}
		t.appendCell(grid, territories, empty)
		*cellCount++
		added = true
		return true
	})
	return added
}

func (t *Territory) setNumDice(n int) {
	if n < 1 {
		n = 1
	}
	if n > MaxDice {
		n = MaxDice
	}
	t.NumDice = n
}
