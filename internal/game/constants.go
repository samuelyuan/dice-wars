package game

const (
	MaxPlayers       = 8
	MinPlayers       = 2
	GridWidth        = 60
	GridHeight       = 40
	NumTerritories   = 80
	TerritorySize    = 25
	MaxDice          = 8
	MaxRemainingDice = 100
	HexRadius        = 10.0

	GrowthInterval    = 30   // milliseconds
	AIStepInterval    = 500  // milliseconds
	AISelectInterval  = 300  // milliseconds
	AIAttackInterval  = 300  // milliseconds
	DiceRollRevealInterval = 1200 // milliseconds — stagger dice one-by-one
	DiceRollHoldInterval   = 800 // milliseconds — brief pause on final totals before resolving
)

// intervalSec converts a millisecond interval to seconds.
func intervalSec(ms int) float64 {
	return float64(ms) / 1000.0
}

// diceRollTotalDuration is reveal + hold before the attack resolves.
func diceRollTotalDuration() float64 {
	return intervalSec(DiceRollRevealInterval) + intervalSec(DiceRollHoldInterval)
}

// Hex directions in axial coordinates (q, r).
var Directions = [6][2]int{{1, 0}, {0, 1}, {-1, 1}, {-1, 0}, {0, -1}, {1, -1}}

// PlayerColors for up to eight players.
var PlayerColors = [MaxPlayers]struct{ R, G, B uint8 }{
	{213, 2, 2},     // red
	{2, 117, 2},     // dark green
	{245, 245, 15},  // yellow
	{2, 230, 230},   // cyan
	{117, 2, 230},   // purple
	{255, 35, 150},  // pink
	{117, 230, 2},   // light green
	{230, 127, 2},   // orange
}
