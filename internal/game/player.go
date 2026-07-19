package game

import "math/rand/v2"

type Player struct {
	Index         int
	Human         bool
	TerritoryIDs  []int
	RemainingDice int
}

// LargestConnectedGroup returns the size of this player's largest connected territory cluster.
func (p *Player) LargestConnectedGroup(territories []*Territory) int {
	visited := make(map[int]bool)
	largest := 0

	for _, startID := range p.TerritoryIDs {
		if !validTerritoryID(territories, startID) || visited[startID] {
			continue
		}
		visited[startID] = true
		size := p.floodFillSize(startID, territories, visited)
		if size > largest {
			largest = size
		}
	}
	return largest
}

func (p *Player) floodFillSize(terrID int, territories []*Territory, visited map[int]bool) int {
	if !validTerritoryID(territories, terrID) {
		return 0
	}

	size := 1
	for _, nbID := range territories[terrID].Neighbours {
		if !validTerritoryID(territories, nbID) {
			continue
		}
		neighbor := territories[nbID]
		if neighbor.Owner != p.Index || visited[nbID] {
			continue
		}
		visited[nbID] = true
		size += p.floodFillSize(nbID, territories, visited)
	}
	return size
}

func (p *Player) addTerritory(terr *Territory, oldOwner *Player) {
	if oldOwner != nil {
		oldOwner.TerritoryIDs = removeInt(oldOwner.TerritoryIDs, terr.ID)
	}
	terr.Owner = p.Index
	p.TerritoryIDs = append(p.TerritoryIDs, terr.ID)
}

func (p *Player) removeTerritory(terr *Territory) {
	p.TerritoryIDs = removeInt(p.TerritoryIDs, terr.ID)
	terr.Owner = -1
}

func (p *Player) addDice(rng *rand.Rand, count int, distribute bool, territories []*Territory) {
	p.RemainingDice += count
	if p.RemainingDice > MaxRemainingDice {
		p.RemainingDice = MaxRemainingDice
	}
	if distribute {
		p.distributeDice(rng, p.RemainingDice, territories)
	}
}

func (p *Player) distributeDice(rng *rand.Rand, count int, territories []*Territory) bool {
	if count > p.RemainingDice {
		count = p.RemainingDice
	}
	territoryCount := len(p.TerritoryIDs)
	if territoryCount == 0 {
		return false
	}

	for placed := 0; placed < count; placed++ {
		if !p.placeOneDie(rng, territories, territoryCount) {
			return false
		}
	}
	return true
}

func (p *Player) placeOneDie(rng *rand.Rand, territories []*Territory, territoryCount int) bool {
	placed := false
	forEachShuffled(territoryCount, rng, func(idx int) bool {
		terr := territories[p.TerritoryIDs[idx]]
		if terr.NumDice >= MaxDice {
			return false
		}
		terr.setNumDice(terr.NumDice + 1)
		p.RemainingDice--
		placed = true
		return true
	})
	return placed
}
