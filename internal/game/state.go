package game

// AttackEffects holds domain outcomes after applying a resolved attack.
type AttackEffects struct {
	GameOver      bool
	VictoryPlayer int
	VictoryHuman  bool
	PlayersLeft   int
}

// ApplyAttackResult applies combat rules for a finished dice roll to territories and players.
func ApplyAttackResult(
	grid *HexGrid,
	territories []*Territory,
	players []*Player,
	attackerTerrID, defenderTerrID int,
	result AttackResult,
	playersLeft int,
) AttackEffects {
	attackerTerr := territories[attackerTerrID]
	defenderTerr := territories[defenderTerrID]
	attacker := players[attackerTerr.Owner]
	defender := players[defenderTerr.Owner]

	if result.AttackerWins {
		attacker.addTerritory(defenderTerr, defender)
		defenderTerr.setNumDice(attackerTerr.NumDice - 1)
		defenderTerr.regenerateNeighbours(grid, territories)
		if len(defender.TerritoryIDs) == 0 {
			playersLeft--
		}
	}

	attackerTerr.setNumDice(1)
	attackerTerr.Selected = false
	defenderTerr.Selected = false

	effects := AttackEffects{
		PlayersLeft:   playersLeft,
		VictoryPlayer: -1,
	}

	if playersLeft == 1 {
		effects.GameOver = true
		effects.VictoryPlayer = attacker.Index
		effects.VictoryHuman = attacker.Human
	}

	return effects
}
