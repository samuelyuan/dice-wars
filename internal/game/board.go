package game

import (
	"math/rand/v2"
	"time"
)

type Phase int

const (
	PhaseIdle Phase = iota
	PhaseDiceRoll
	PhaseGrowing
	PhaseAIWait
	PhaseAISelect
	PhaseAIAttack
)

type Board struct {
	Grid          *HexGrid
	Territories   []*Territory
	Players       []*Player
	NumPlayers    int
	PlayerTurn    int
	PlayersLeft   int
	SelectedTerr  int // territory ID, -1 if none
	OtherTerr     int
	Phase         Phase
	CheatMode     bool
	AutoMode      bool
	StatusMessage string
	LastAttack    AttackResult
	VictoryPlayer int
	VictoryHuman  bool
	GameOver      bool
	phaseTimer    float64
	growPending   bool
	seed          uint64
	rng           *rand.Rand
}

const minTerritoryCells = 6

func NewBoard(numPlayers int, humanList []bool) *Board {
	return NewBoardWithSeed(numPlayers, humanList, uint64(time.Now().UnixNano()))
}

func NewBoardWithSeed(numPlayers int, humanList []bool, seed uint64) *Board {
	b := &Board{
		NumPlayers:    numPlayers,
		SelectedTerr:  -1,
		OtherTerr:     -1,
		VictoryPlayer: -1,
		seed:          seed,
	}
	b.setup(humanList)
	return b
}

func (b *Board) Initialize() {
	humanList := make([]bool, b.NumPlayers)
	for i := range humanList {
		if i < len(b.Players) {
			humanList[i] = b.Players[i].Human
		}
	}
	b.setup(humanList)
}

func (b *Board) setup(humanList []bool) {
	b.rng = rngFromSeed(b.seed)
	b.resetForNewGame()
	b.createPlayers(humanList)
	b.generateTerritories()
	b.removeSmallTerritories()
	b.dealInitialDice()
	b.beginTurn(b.rng.IntN(b.NumPlayers))
}

func (b *Board) resetForNewGame() {
	b.Grid = NewHexGrid()
	b.Territories = nil
	b.Players = nil
	b.SelectedTerr = -1
	b.OtherTerr = -1
	b.Phase = PhaseIdle
	b.GameOver = false
	b.VictoryPlayer = -1
	b.AutoMode = false
	b.StatusMessage = ""
	b.phaseTimer = 0
	b.growPending = false
}

func (b *Board) createPlayers(humanList []bool) {
	b.PlayersLeft = b.NumPlayers
	b.Players = make([]*Player, 0, b.NumPlayers)
	for i := 0; i < b.NumPlayers; i++ {
		b.Players = append(b.Players, &Player{Index: i, Human: humanList[i]})
	}
}

func (b *Board) removeSmallTerritories() {
	var kept []*Territory
	for _, t := range b.Territories {
		if len(t.CellIDs) >= minTerritoryCells {
			kept = append(kept, t)
			continue
		}
		for _, ax := range t.CellIDs {
			if h := b.Grid.Hexes[ax]; h != nil {
				h.TerrID = -1
			}
		}
		if t.Owner >= 0 {
			b.Players[t.Owner].removeTerritory(t)
		}
	}
	b.Territories = kept
	b.reindexTerritories()
}

func (b *Board) dealInitialDice() {
	initialDice := len(b.Territories) * 15 / 10 / b.NumPlayers
	for _, p := range b.Players {
		p.addDice(b.rng, initialDice, true, b.Territories)
	}
}

func (b *Board) reindexTerritories() {
	for i, t := range b.Territories {
		t.ID = i
		for _, ax := range t.CellIDs {
			b.Grid.Hexes[ax].TerrID = i
		}
	}

	for _, p := range b.Players {
		p.TerritoryIDs = p.TerritoryIDs[:0]
	}
	for _, t := range b.Territories {
		if t.Owner >= 0 {
			b.Players[t.Owner].TerritoryIDs = append(b.Players[t.Owner].TerritoryIDs, t.ID)
		}
	}

	for _, t := range b.Territories {
		t.regenerateNeighbours(b.Grid, b.Territories)
	}
}

func (b *Board) generateTerritories() {
	x := b.rng.IntN(GridWidth)
	y := b.rng.IntN(GridHeight)
	ax := Axial{Q: x - (y - (y & 1)) / 2, R: y}
	hex := b.Grid.Hexes[ax]

	terr := b.newTerritory(0)
	terr.appendCell(b.Grid, b.Territories, hex)
	terr.grow(b.rng, b.Grid, b.Territories, TerritorySize)

	for terrCount := 1; terrCount < NumTerritories; terrCount++ {
		seedHex := terr.findEmptyAdjacent(b.rng, b.Grid)
		for seedHex == nil {
			seedHex = b.Territories[b.rng.IntN(terrCount)].findEmptyAdjacent(b.rng, b.Grid)
		}
		terr = b.newTerritory(terrCount % b.NumPlayers)
		terr.appendCell(b.Grid, b.Territories, seedHex)
		terr.grow(b.rng, b.Grid, b.Territories, TerritorySize)
	}
}

func (b *Board) newTerritory(ownerIdx int) *Territory {
	t := &Territory{
		ID:      len(b.Territories),
		Owner:   ownerIdx,
		NumDice: 1,
	}
	b.Territories = append(b.Territories, t)
	b.Players[ownerIdx].TerritoryIDs = append(b.Players[ownerIdx].TerritoryIDs, t.ID)
	return t
}

func (b *Board) IsBusy() bool {
	return b.Phase != PhaseIdle
}

func (b *Board) Click(px, py float64) {
	if b.GameOver || b.IsBusy() {
		return
	}
	hex := b.Grid.PickHex(px, py)
	if hex == nil || hex.TerrID < 0 {
		return
	}
	clicked := b.Territories[hex.TerrID]
	if clicked.Owner < 0 {
		return
	}
	if b.tryToggleSelection(clicked) {
		return
	}
	if b.tryAttack(clicked) {
		return
	}
	if !b.canSelectTerritory(clicked) {
		return
	}
	b.selectTerritory(clicked)
}

func (b *Board) tryToggleSelection(terr *Territory) bool {
	if b.SelectedTerr < 0 || b.SelectedTerr != terr.ID {
		return false
	}
	b.clearSelection()
	return true
}

func (b *Board) tryAttack(target *Territory) bool {
	if b.SelectedTerr < 0 {
		return false
	}
	attacker := b.Territories[b.SelectedTerr]
	if attacker.Owner == target.Owner || !containsInt(target.Neighbours, attacker.ID) {
		return false
	}
	b.OtherTerr = target.ID
	b.processAttack()
	return true
}

func (b *Board) canSelectTerritory(terr *Territory) bool {
	return terr.Owner == b.PlayerTurn && terr.NumDice >= 2
}

func (b *Board) selectTerritory(terr *Territory) {
	b.clearSelection()
	terr.Selected = true
	b.SelectedTerr = terr.ID
}

func (b *Board) clearSelection() {
	if b.SelectedTerr >= 0 {
		b.Territories[b.SelectedTerr].Selected = false
	}
	b.SelectedTerr = -1
	b.OtherTerr = -1
}

func (b *Board) processAttack() {
	attackerTerr := b.Territories[b.SelectedTerr]
	defenderTerr := b.Territories[b.OtherTerr]
	attackerTerr.Selected = true
	defenderTerr.Selected = true

	attacker := b.Players[attackerTerr.Owner]
	defender := b.Players[defenderTerr.Owner]
	cheatAttacker := b.CheatMode && attacker.Human
	cheatDefender := b.CheatMode && defender.Human

	b.LastAttack = ResolveAttack(b.rng, attackerTerr.NumDice, defenderTerr.NumDice, cheatAttacker, cheatDefender)
	b.Phase = PhaseDiceRoll
	b.phaseTimer = 0
}

func (b *Board) attackFinished() {
	attackerIdx := b.Territories[b.SelectedTerr].Owner
	effects := ApplyAttackResult(
		b.Grid, b.Territories, b.Players,
		b.SelectedTerr, b.OtherTerr,
		b.LastAttack, b.PlayersLeft,
	)
	b.PlayersLeft = effects.PlayersLeft
	b.clearSelection()

	if effects.GameOver {
		b.GameOver = true
		b.VictoryPlayer = effects.VictoryPlayer
		b.VictoryHuman = effects.VictoryHuman
		b.Phase = PhaseIdle
		return
	}

	b.resumeAfterAttack(attackerIdx)
}

func (b *Board) resumeAfterAttack(attackerIdx int) {
	if b.AutoMode || !b.Players[attackerIdx].Human {
		b.Phase = PhaseAIWait
		b.phaseTimer = intervalSec(AIStepInterval)
		return
	}
	b.Phase = PhaseIdle
}

func (b *Board) EndTurn() {
	if b.GameOver || b.IsBusy() {
		return
	}
	b.concludeTurnWithReinforcements()
}

func (b *Board) concludeTurnWithReinforcements() {
	if b.GameOver {
		return
	}
	b.AutoMode = false
	player := b.Players[b.PlayerTurn]
	player.addDice(b.rng, player.LargestConnectedGroup(b.Territories), false, b.Territories)
	b.Phase = PhaseGrowing
	b.phaseTimer = intervalSec(GrowthInterval)
	b.growPending = true
}

func (b *Board) growStep() {
	player := b.Players[b.PlayerTurn]
	if player.distributeDice(b.rng, 1, b.Territories) && player.RemainingDice > 0 {
		b.phaseTimer = intervalSec(GrowthInterval)
		return
	}
	b.advancePlayerTurn()
}

func (b *Board) advancePlayerTurn() {
	b.growPending = false
	b.Phase = PhaseIdle
	b.advanceToNextActivePlayer()
	b.clearSelection()
	b.AutoMode = false
	b.updateStatusForTurn()

	if !b.Players[b.PlayerTurn].Human {
		b.scheduleAITurn()
	}
}

func (b *Board) advanceToNextActivePlayer() {
	for {
		b.PlayerTurn = (b.PlayerTurn + 1) % b.NumPlayers
		if len(b.Players[b.PlayerTurn].TerritoryIDs) > 0 {
			return
		}
	}
}

func (b *Board) beginTurn(playerIdx int) {
	b.PlayerTurn = playerIdx
	b.updateStatusForTurn()
	if !b.Players[b.PlayerTurn].Human {
		b.scheduleAITurn()
	}
}

func (b *Board) updateStatusForTurn() {
	if b.Players[b.PlayerTurn].Human && !b.AutoMode {
		b.StatusMessage = "Your turn!"
		return
	}
	b.StatusMessage = ""
}

func (b *Board) StartAITurn() {
	if b.GameOver || b.IsBusy() {
		return
	}
	if b.Players[b.PlayerTurn].Human {
		b.AutoMode = true
	}
	b.scheduleAITurn()
}

func (b *Board) scheduleAITurn() {
	b.Phase = PhaseAIWait
	b.phaseTimer = intervalSec(AIStepInterval)
}

func (b *Board) nextAIStep() {
	player := b.Players[b.PlayerTurn]
	if len(player.TerritoryIDs) == 0 {
		b.advancePlayerTurn()
		return
	}

	attackerID, defenderID, found := b.findAIAttackTarget(player)
	if found {
		b.SelectedTerr = attackerID
		b.OtherTerr = defenderID
		b.Phase = PhaseAISelect
		b.phaseTimer = 0
		return
	}

	b.concludeTurnWithReinforcements()
}

func (b *Board) findAIAttackTarget(player *Player) (attackerID, defenderID int, found bool) {
	forEachShuffled(len(player.TerritoryIDs), b.rng, func(terrIdx int) bool {
		tid := player.TerritoryIDs[terrIdx]
		attacker := b.Territories[tid]
		if attacker.NumDice < 2 || len(attacker.Neighbours) == 0 {
			return false
		}

		maxDice := maxAttackableDice(attacker.NumDice)
		forEachShuffled(len(attacker.Neighbours), b.rng, func(nbIdx int) bool {
			nbID := attacker.Neighbours[nbIdx]
			defender := b.Territories[nbID]
			if defender.Owner < 0 || defender.Owner == player.Index {
				return false
			}
			if defender.NumDice > maxDice {
				return false
			}
			attackerID, defenderID, found = tid, nbID, true
			return true
		})
		return found
	})
	return
}

func maxAttackableDice(numDice int) int {
	if numDice >= MaxDice {
		return numDice
	}
	return numDice - 1
}

func (b *Board) aiSelectStep() {
	attacker := b.Territories[b.SelectedTerr]
	if !attacker.Selected {
		attacker.Selected = true
		b.phaseTimer = intervalSec(AISelectInterval)
		return
	}
	b.Territories[b.OtherTerr].Selected = true
	b.Phase = PhaseAIAttack
	b.phaseTimer = intervalSec(AIAttackInterval)
}

func (b *Board) aiAttackStep() {
	b.processAttack()
}

func (b *Board) Update(dt float64) {
	if b.GameOver {
		return
	}

	switch b.Phase {
	case PhaseDiceRoll:
		b.phaseTimer += dt
		if b.phaseTimer >= diceRollTotalDuration() {
			b.attackFinished()
		}
	case PhaseGrowing:
		b.tickCountdown(dt, func() {
			if b.growPending {
				b.growStep()
			}
		})
	case PhaseAIWait:
		b.tickCountdown(dt, b.nextAIStep)
	case PhaseAISelect:
		b.tickCountdown(dt, b.aiSelectStep)
	case PhaseAIAttack:
		b.tickCountdown(dt, b.aiAttackStep)
	}
}

func (b *Board) tickCountdown(dt float64, onExpire func()) {
	b.phaseTimer -= dt
	if b.phaseTimer <= 0 {
		onExpire()
	}
}

func (b *Board) ConnectedTerrCount(playerIdx int) int {
	if !inRange(playerIdx, len(b.Players)) {
		return 0
	}
	return b.Players[playerIdx].LargestConnectedGroup(b.Territories)
}

func (b *Board) IsPlayerActive(playerIdx int) bool {
	if !inRange(playerIdx, len(b.Players)) {
		return false
	}
	return len(b.Players[playerIdx].TerritoryIDs) > 0
}

// HumanIndex returns the index of the human player, or -1 if there isn't one.
func (b *Board) HumanIndex() int {
	for i, p := range b.Players {
		if p.Human {
			return i
		}
	}
	return -1
}

// HumanEliminated reports whether the human player has lost all territory
// but the game hasn't concluded yet (other players are still fighting it
// out). Returns false once GameOver is true — that's a normal end-of-game,
// not a mid-game elimination.
func (b *Board) HumanEliminated() bool {
	if b.GameOver {
		return false
	}
	idx := b.HumanIndex()
	return idx >= 0 && !b.IsPlayerActive(idx)
}
