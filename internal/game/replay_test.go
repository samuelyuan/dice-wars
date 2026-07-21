package game

import "testing"

// findAttackFor returns an attacker/defender pair using the AI's targeting
// rules, or found=false if none exists. Used to script human turns in tests.
func findAttackFor(board *Board, playerIdx int) (attackerID, defenderID int, found bool) {
	player := board.Players[playerIdx]
	for _, tid := range player.TerritoryIDs {
		attacker := board.Territories[tid]
		if attacker.NumDice < 2 || len(attacker.Neighbours) == 0 {
			continue
		}
		maxDice := maxAttackableDice(attacker.NumDice)
		for _, nbID := range attacker.Neighbours {
			defender := board.Territories[nbID]
			if defender.Owner < 0 || defender.Owner == playerIdx {
				continue
			}
			if defender.NumDice > maxDice {
				continue
			}
			return tid, nbID, true
		}
	}
	return 0, 0, false
}

// driveHumanTurn plays player 0's turn via the same internal calls Click()
// makes, skipping pixel-to-hex resolution — deriving coordinates from cell
// centers proved unreliable on odd territory shapes.
func driveHumanTurn(board *Board, dt float64) {
	for board.PlayerTurn == 0 && !board.GameOver {
		if board.IsBusy() {
			board.Update(dt)
			continue
		}
		attackerID, defenderID, found := findAttackFor(board, 0)
		if !found {
			board.EndTurn()
			board.Update(dt)
			continue
		}
		board.selectTerritory(board.Territories[attackerID])
		board.OtherTerr = defenderID
		board.processAttack()
		for board.IsBusy() && !board.GameOver {
			board.Update(dt)
		}
	}
}

// boardSnapshot captures per-territory owner/dice for full-state comparison.
type boardSnapshot struct {
	owners []int
	dice   []int
}

func snapshotBoard(b *Board) boardSnapshot {
	s := boardSnapshot{owners: make([]int, len(b.Territories)), dice: make([]int, len(b.Territories))}
	for i, t := range b.Territories {
		s.owners[i] = t.Owner
		s.dice[i] = t.NumDice
	}
	return s
}

func (s boardSnapshot) diff(other boardSnapshot) (mismatches int) {
	for i := range s.owners {
		if s.owners[i] != other.owners[i] || s.dice[i] != other.dice[i] {
			mismatches++
		}
	}
	return mismatches
}

// TestReplayReproducesOriginalGame is the regression guard for record-and-
// replay: plays a mixed human/AI game, exports and replays it via Board.Update
// alone, and asserts the final state matches exactly.
func TestReplayReproducesOriginalGame(t *testing.T) {
	const dt = 1.0 / 60.0
	const maxIterations = 500_000

	seeds := []uint64{1, 2, 3, 42, 100}
	for _, seed := range seeds {
		humanList := []bool{true, false, false, false, false, false}
		board := NewBoardWithSeed(6, humanList, seed)

		iterations := 0
		for !board.GameOver && iterations < maxIterations {
			if board.PlayerTurn == 0 {
				driveHumanTurn(board, dt)
			} else {
				board.Update(dt)
			}
			iterations++
		}
		if iterations >= maxIterations {
			t.Fatalf("seed %d: game didn't conclude within %d iterations", seed, maxIterations)
		}
		if len(board.Actions) == 0 {
			t.Fatalf("seed %d: no actions recorded", seed)
		}

		replay := board.ExportReplay()
		rp := NewReplayPlayer(replay)
		steps := 0
		const maxSteps = 500_000
		for !rp.Finished() && steps < maxSteps {
			rp.Board.Update(dt)
			steps++
		}
		if steps >= maxSteps {
			t.Fatalf("seed %d: replay never converged (stuck) after %d steps", seed, steps)
		}

		if !rp.Board.GameOver {
			t.Errorf("seed %d: replay didn't reach GameOver (exhausted early?) actionsPlayed=%d/%d",
				seed, rp.ActionsPlayed(), len(replay.Actions))
		}
		if rp.Board.VictoryPlayer != board.VictoryPlayer {
			t.Errorf("seed %d: VictoryPlayer mismatch: original=%d replay=%d", seed, board.VictoryPlayer, rp.Board.VictoryPlayer)
		}

		orig := snapshotBoard(board)
		rep := snapshotBoard(rp.Board)
		if mismatches := orig.diff(rep); mismatches > 0 {
			t.Errorf("seed %d: %d/%d territories mismatched between original and replay",
				seed, mismatches, len(board.Territories))
		}

		for i, p := range board.Players {
			if len(p.TerritoryIDs) != len(rp.Board.Players[i].TerritoryIDs) {
				t.Errorf("seed %d: player %d territory count mismatch: original=%d replay=%d",
					seed, i, len(p.TerritoryIDs), len(rp.Board.Players[i].TerritoryIDs))
			}
		}
	}
}

// TestReplayFeedPeekMatchingDoesNotAdvance verifies peekMatching never moves
// the cursor, since nextAIStep depends on looking ahead without consuming.
func TestReplayFeedPeekMatchingDoesNotAdvance(t *testing.T) {
	f := &replayFeed{actions: []Action{
		{Type: ActionAttack, PlayerIndex: 1, TerritoryID: 5, OtherTerrID: 9},
	}}

	if action, ok := f.peekMatching(ActionAttack, 1); !ok || action.TerritoryID != 5 {
		t.Fatalf("peekMatching on match: got action=%+v ok=%v, want TerritoryID=5 ok=true", action, ok)
	}
	if f.index != 0 {
		t.Fatalf("peekMatching on match advanced the cursor to %d, want 0", f.index)
	}

	if _, ok := f.peekMatching(ActionTurnEnd, 1); ok {
		t.Fatal("peekMatching matched on wrong action type")
	}
	if _, ok := f.peekMatching(ActionAttack, 0); ok {
		t.Fatal("peekMatching matched on wrong player index")
	}
	if f.index != 0 {
		t.Fatalf("peekMatching on mismatch advanced the cursor to %d, want 0", f.index)
	}
}

// TestReplayFeedConsumeMatching verifies the cursor only advances on a match.
func TestReplayFeedConsumeMatching(t *testing.T) {
	f := &replayFeed{actions: []Action{
		{Type: ActionGrowthPlace, PlayerIndex: 2, TerritoryID: 3},
		{Type: ActionTurnEnd, PlayerIndex: 2},
	}}

	if _, ok := f.consumeMatching(ActionGrowthPlace, 0); ok {
		t.Fatal("consumeMatching matched on wrong player index")
	}
	if f.index != 0 {
		t.Fatalf("consumeMatching on mismatch advanced the cursor to %d, want 0", f.index)
	}

	action, ok := f.consumeMatching(ActionGrowthPlace, 2)
	if !ok || action.TerritoryID != 3 {
		t.Fatalf("consumeMatching on match: got action=%+v ok=%v, want TerritoryID=3 ok=true", action, ok)
	}
	if f.index != 1 {
		t.Fatalf("consumeMatching on match left cursor at %d, want 1", f.index)
	}

	if _, ok := f.consumeMatching(ActionTurnEnd, 2); !ok {
		t.Fatal("consumeMatching failed to match the second action")
	}
	if !f.exhausted() {
		t.Fatal("feed should be exhausted after consuming both actions")
	}
}

// TestReplayFeedPeekThenConsumeSameAction mirrors the nextAIStep ->
// processAttack split: peek to decide, consume once acted on. This broke
// when both calls advanced the cursor.
func TestReplayFeedPeekThenConsumeSameAction(t *testing.T) {
	f := &replayFeed{actions: []Action{
		{Type: ActionAttack, PlayerIndex: 0, TerritoryID: 7, OtherTerrID: 8},
	}}

	peeked, ok := f.peekMatching(ActionAttack, 0)
	if !ok {
		t.Fatal("peekMatching didn't find the action")
	}

	consumed, ok := f.consumeMatching(ActionAttack, 0)
	if !ok {
		t.Fatal("consumeMatching didn't find the same action after peeking")
	}
	if consumed.TerritoryID != peeked.TerritoryID || consumed.OtherTerrID != peeked.OtherTerrID {
		t.Fatalf("consumed action %+v differs from peeked action %+v", consumed, peeked)
	}
	if !f.exhausted() {
		t.Fatal("feed should be exhausted after the single consume")
	}
}

// TestReplaySeekLandsAtExactAction verifies SeekToProgress lands exactly and never hangs.
func TestReplaySeekLandsAtExactAction(t *testing.T) {
	const dt = 1.0 / 60.0
	humanList := []bool{true, false, false, false, false, false}
	board := NewBoardWithSeed(6, humanList, 7)

	iterations := 0
	for !board.GameOver && iterations < 500_000 {
		if board.PlayerTurn == 0 {
			driveHumanTurn(board, dt)
		} else {
			board.Update(dt)
		}
		iterations++
	}
	if !board.GameOver {
		t.Fatal("game didn't conclude")
	}

	replay := board.ExportReplay()
	if len(replay.Actions) == 0 {
		t.Fatal("no actions recorded")
	}

	for _, f := range []float64{0, 0.1, 0.25, 0.5, 0.75, 0.9, 1.0} {
		rp := NewReplayPlayer(replay)
		rp.SeekToProgress(f)
		want := int(f * float64(len(replay.Actions)))
		if got := rp.ActionsPlayed(); got != want && !rp.Finished() {
			t.Errorf("seek to %.2f: actionsPlayed=%d, want %d (Finished=%v)", f, got, want, rp.Finished())
		}
	}
}
