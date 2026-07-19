package game

import "testing"

func TestNewBoardWithSeedDeterministic(t *testing.T) {
	const seed uint64 = 42
	humans := []bool{true, false, false}

	a := NewBoardWithSeed(3, humans, seed)
	b := NewBoardWithSeed(3, humans, seed)

	if len(a.Territories) != len(b.Territories) {
		t.Fatalf("territory count: %d vs %d", len(a.Territories), len(b.Territories))
	}
	for i := range a.Territories {
		ta, tb := a.Territories[i], b.Territories[i]
		if ta.Owner != tb.Owner || len(ta.CellIDs) != len(tb.CellIDs) {
			t.Fatalf("territory %d differs between runs", i)
		}
		for j, ax := range ta.CellIDs {
			if ax != tb.CellIDs[j] {
				t.Fatalf("territory %d cell %d differs", i, j)
			}
		}
	}
	if a.PlayerTurn != b.PlayerTurn {
		t.Fatalf("player turn: %d vs %d", a.PlayerTurn, b.PlayerTurn)
	}
	for i, pa := range a.Players {
		pb := b.Players[i]
		if len(pa.TerritoryIDs) != len(pb.TerritoryIDs) || pa.RemainingDice != pb.RemainingDice {
			t.Fatalf("player %d state differs", i)
		}
	}
}

func TestInitializeReproducibleWithSeed(t *testing.T) {
	const seed uint64 = 99
	humans := []bool{true, false}
	b := NewBoardWithSeed(2, humans, seed)
	firstTerrCount := len(b.Territories)
	firstTurn := b.PlayerTurn

	b.Initialize()

	if len(b.Territories) != firstTerrCount {
		t.Fatalf("restart changed territory count: %d -> %d", firstTerrCount, len(b.Territories))
	}
	if b.PlayerTurn != firstTurn {
		t.Fatalf("restart changed first player: %d -> %d", firstTurn, b.PlayerTurn)
	}
}

func TestApplyAttackResultAttackerWins(t *testing.T) {
	grid := NewHexGrid()
	players := []*Player{
		{Index: 0, Human: true},
		{Index: 1, Human: false},
	}
	territories := []*Territory{
		{ID: 0, Owner: 0, NumDice: 4, Neighbours: []int{1}},
		{ID: 1, Owner: 1, NumDice: 2, Neighbours: []int{0}},
	}
	players[0].TerritoryIDs = []int{0}
	players[1].TerritoryIDs = []int{1}

	result := AttackResult{
		AttackTotal:  18,
		DefenseTotal: 7,
		AttackerWins: true,
	}
	effects := ApplyAttackResult(grid, territories, players, 0, 1, result, 2)

	if territories[1].Owner != 0 {
		t.Fatalf("defender owner want 0 got %d", territories[1].Owner)
	}
	if territories[0].NumDice != 1 {
		t.Fatalf("attacker dice want 1 got %d", territories[0].NumDice)
	}
	if territories[1].NumDice != 3 {
		t.Fatalf("conquered dice want 3 got %d", territories[1].NumDice)
	}
	if len(players[1].TerritoryIDs) != 0 {
		t.Fatalf("defender should have no territories")
	}
	if effects.PlayersLeft != 1 {
		t.Fatalf("players left want 1 got %d", effects.PlayersLeft)
	}
	if !effects.GameOver || effects.VictoryPlayer != 0 || !effects.VictoryHuman {
		t.Fatalf("expected player 0 human victory, got %+v", effects)
	}
}

func TestApplyAttackResultDefenderWins(t *testing.T) {
	grid := NewHexGrid()
	players := []*Player{
		{Index: 0},
		{Index: 1},
	}
	territories := []*Territory{
		{ID: 0, Owner: 0, NumDice: 3, Neighbours: []int{1}},
		{ID: 1, Owner: 1, NumDice: 4, Neighbours: []int{0}},
	}
	players[0].TerritoryIDs = []int{0}
	players[1].TerritoryIDs = []int{1}

	result := AttackResult{AttackTotal: 5, DefenseTotal: 12, AttackerWins: false}
	effects := ApplyAttackResult(grid, territories, players, 0, 1, result, 2)

	if territories[1].Owner != 1 {
		t.Fatalf("defender should keep territory")
	}
	if territories[0].NumDice != 1 {
		t.Fatalf("attacker dice want 1 got %d", territories[0].NumDice)
	}
	if effects.GameOver {
		t.Fatal("game should not be over")
	}
	if effects.PlayersLeft != 2 {
		t.Fatalf("players left want 2 got %d", effects.PlayersLeft)
	}
}

func TestLargestConnectedGroup(t *testing.T) {
	territories := []*Territory{
		{ID: 0, Owner: 0, Neighbours: []int{1}},
		{ID: 1, Owner: 0, Neighbours: []int{0, 2}},
		{ID: 2, Owner: 0, Neighbours: []int{1}},
		{ID: 3, Owner: 0, Neighbours: nil},
	}
	p := &Player{Index: 0, TerritoryIDs: []int{0, 1, 2, 3}}

	if got := p.LargestConnectedGroup(territories); got != 3 {
		t.Fatalf("largest connected group want 3 got %d", got)
	}
}
