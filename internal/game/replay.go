package game

// ActionType identifies what a recorded event represents.
type ActionType int

const (
	// ActionAttack carries resolved dice rolls, reused during replay.
	ActionAttack ActionType = iota
	// ActionGrowthPlace is one reinforcement die placed on a territory.
	ActionGrowthPlace
	// ActionTurnEnd marks a player's turn concluding.
	ActionTurnEnd
)

// Action is a single recorded event, in chronological order across all
// players. Replaying the sequence exactly reproduces the game.
type Action struct {
	Type          ActionType `json:"type"`
	PlayerIndex   int        `json:"player_index"`
	TerritoryID   int        `json:"territory_id"`             // attacker (ActionAttack) or placement target (ActionGrowthPlace)
	OtherTerrID   int        `json:"other_terr_id"`             // defender, ActionAttack only
	AttackerRolls []int      `json:"attacker_rolls,omitempty"` // ActionAttack only
	DefenderRolls []int      `json:"defender_rolls,omitempty"` // ActionAttack only
}

// Replay represents a complete game record that can be replayed
type Replay struct {
	Seed          uint64   `json:"seed"`
	NumPlayers    int      `json:"num_players"`
	HumanList     []bool   `json:"human_list"`
	Actions       []Action `json:"actions"`
	VictoryPlayer int      `json:"victory_player"` // -1 if exported before the game concluded
	VictoryHuman  bool     `json:"victory_human"`
}

// IsPartialGame reports whether this replay was exported mid-game.
func (r *Replay) IsPartialGame() bool {
	return r.VictoryPlayer < 0
}

// RecordAction appends an event to the board's action history.
func (b *Board) RecordAction(a Action) {
	if b.replay != nil {
		return // this board is itself replaying; don't re-record what it plays back
	}
	b.Actions = append(b.Actions, a)
}

// ExportReplay creates a Replay object from the current board state
func (b *Board) ExportReplay() *Replay {
	humanList := make([]bool, b.NumPlayers)
	for i, p := range b.Players {
		humanList[i] = p.Human
	}
	return &Replay{
		Seed:          b.seed,
		NumPlayers:    b.NumPlayers,
		HumanList:     humanList,
		Actions:       append([]Action{}, b.Actions...), // Copy actions
		VictoryPlayer: b.VictoryPlayer,
		VictoryHuman:  b.VictoryHuman,
	}
}

// replayFeed is a read cursor over a recorded action log.
type replayFeed struct {
	actions []Action
	index   int
}

// peek returns the next unconsumed action without advancing the cursor.
func (f *replayFeed) peek() (Action, bool) {
	if f.index >= len(f.actions) {
		return Action{}, false
	}
	return f.actions[f.index], true
}

// advance consumes the next action.
func (f *replayFeed) advance() {
	f.index++
}

// peekMatching returns the next action without consuming it, if it matches
// type and player — for looking ahead before acting (see nextAIStep).
func (f *replayFeed) peekMatching(actionType ActionType, playerIdx int) (Action, bool) {
	action, ok := f.peek()
	if !ok || action.Type != actionType || action.PlayerIndex != playerIdx {
		return Action{}, false
	}
	return action, true
}

// consumeMatching consumes the next action if it matches type and player index.
func (f *replayFeed) consumeMatching(actionType ActionType, playerIdx int) (Action, bool) {
	action, ok := f.peekMatching(actionType, playerIdx)
	if !ok {
		return Action{}, false
	}
	f.advance()
	return action, true
}

// exhausted reports whether every recorded action has been consumed.
func (f *replayFeed) exhausted() bool {
	return f.index >= len(f.actions)
}

// replayStepDT is the fixed timestep for playback and seek fast-forwarding.
const replayStepDT = 1.0 / 60.0

// replaySeekIterationCap guards SeekTo against spinning on a malformed replay.
const replaySeekIterationCap = 5_000_000

// ReplayPlayer drives a Board from a recorded replay via Board.Update.
type ReplayPlayer struct {
	Board  *Board
	replay *Replay
}

// NewReplayPlayer creates a replay-driven board from the given recording.
func NewReplayPlayer(r *Replay) *ReplayPlayer {
	rp := &ReplayPlayer{replay: r}
	rp.reset()
	return rp
}

func (rp *ReplayPlayer) reset() {
	rp.Board = newReplayBoard(rp.replay)
}

// Finished reports whether the game concluded or all actions are consumed.
func (rp *ReplayPlayer) Finished() bool {
	// Avoids reporting finished while the last dice-roll animation plays.
	return rp.Board.GameOver || (rp.Board.ReplayExhausted() && !rp.Board.IsBusy())
}

// ActionsPlayed returns how many actions have been consumed (e.g. "move 42/68").
func (rp *ReplayPlayer) ActionsPlayed() int {
	return rp.Board.replay.index
}

// Progress returns playback position as a fraction in [0,1], for the seek bar.
func (rp *ReplayPlayer) Progress() float64 {
	if rp.Finished() {
		return 1
	}
	total := len(rp.replay.Actions)
	if total == 0 {
		return 0
	}
	return float64(rp.ActionsPlayed()) / float64(total)
}

// SeekToProgress jumps to roughly the given fraction in [0,1] by rebuilding
// the board and fast-forwarding silently (a full game simulates in ms).
func (rp *ReplayPlayer) SeekToProgress(fraction float64) {
	if fraction < 0 {
		fraction = 0
	}
	if fraction > 1 {
		fraction = 1
	}
	targetIndex := int(fraction * float64(len(rp.replay.Actions)))

	rp.reset()
	for i := 0; i < replaySeekIterationCap && !rp.Finished() && rp.ActionsPlayed() < targetIndex; i++ {
		rp.Board.Update(replayStepDT)
	}
}
