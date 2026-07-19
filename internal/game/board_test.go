package game

import "testing"

func TestBoardInitialize(t *testing.T) {
	for players := MinPlayers; players <= MaxPlayers; players++ {
		humans := make([]bool, players)
		humans[0] = true
		for trial := 0; trial < 20; trial++ {
			b := NewBoard(players, humans)
			if len(b.Territories) == 0 {
				t.Fatalf("players=%d trial=%d: no territories", players, trial)
			}
			for i, terr := range b.Territories {
				if terr.ID != i {
					t.Fatalf("players=%d trial=%d: territory %d has id %d", players, trial, i, terr.ID)
				}
				for _, ax := range terr.CellIDs {
					h := b.Grid.Hexes[ax]
					if h.TerrID != i {
						t.Fatalf("players=%d trial=%d: hex terr mismatch want %d got %d", players, trial, i, h.TerrID)
					}
				}
				for _, nbID := range terr.Neighbours {
					if nbID < 0 || nbID >= len(b.Territories) {
						t.Fatalf("players=%d trial=%d: invalid neighbour %d for territory %d", players, trial, nbID, i)
					}
				}
			}
			for _, p := range b.Players {
				conn := p.LargestConnectedGroup(b.Territories)
				if conn > len(p.TerritoryIDs) {
					t.Fatalf("players=%d trial=%d: connected %d > territories %d",
						players, trial, conn, len(p.TerritoryIDs))
				}
				if len(p.TerritoryIDs) > 0 && conn < 1 {
					t.Fatalf("players=%d trial=%d: player has territories but conn=0", players, trial)
				}
			}
		}
	}
}

func TestAIEndTurnGrantsReinforcements(t *testing.T) {
	humans := []bool{false, false, false}
	b := NewBoard(3, humans)
	b.PlayerTurn = 0
	p := b.Players[0]
	before := p.RemainingDice
	b.concludeTurnWithReinforcements()
	if p.RemainingDice <= before {
		conn := p.LargestConnectedGroup(b.Territories)
		t.Fatalf("expected reinforcements, before=%d after=%d conn=%d", before, p.RemainingDice, conn)
	}
	if b.Phase != PhaseGrowing {
		t.Fatalf("expected growing phase, got %v", b.Phase)
	}
}

func TestAITimerFires(t *testing.T) {
	humans := []bool{false, false, false}
	b := NewBoard(3, humans)
	if b.Phase != PhaseAIWait {
		t.Fatalf("expected AI wait at start, got %v", b.Phase)
	}
	// One second exceeds the 500ms AI step interval.
	b.Update(1.0)
	if b.Phase == PhaseAIWait {
		t.Fatal("AI still waiting after 1s — timer may be broken")
	}
}
