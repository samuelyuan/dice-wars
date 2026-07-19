package game

import "strconv"

// IsHumanTurn reports whether the active player has manual control.
func (b *Board) IsHumanTurn() bool {
	if b.GameOver {
		return false
	}
	return b.Players[b.PlayerTurn].Human && !b.AutoMode
}

// TurnBanner returns the active player label shown at the top of the board.
func (b *Board) TurnBanner() string {
	if b.GameOver {
		return "Game over"
	}
	return playerLabel(b.Players[b.PlayerTurn])
}

func playerLabel(p *Player) string {
	n := strconv.Itoa(p.Index + 1)
	if p.Human {
		return "Player " + n + " (you)"
	}
	return "Player " + n + " (CPU)"
}

// AttackAttackerOwner returns the owner index during an attack animation, or -1.
func (b *Board) AttackAttackerOwner() int {
	return b.territoryOwner(b.SelectedTerr)
}

// AttackDefenderOwner returns the defender owner index during an attack animation, or -1.
func (b *Board) AttackDefenderOwner() int {
	return b.territoryOwner(b.OtherTerr)
}

func (b *Board) territoryOwner(terrID int) int {
	if terrID < 0 {
		return -1
	}
	return b.Territories[terrID].Owner
}

// RollRevealProgress returns 0..1 progress through the battle dice animation.
func (b *Board) RollRevealProgress() float64 {
	if b.Phase != PhaseDiceRoll {
		return 1
	}
	dur := intervalSec(DiceRollRevealInterval)
	if dur <= 0 {
		return 1
	}
	progress := b.phaseTimer / dur
	if progress > 1 {
		return 1
	}
	return progress
}

// RevealedAttackerDice returns how many attacker dice are visible during the roll animation.
func (b *Board) RevealedAttackerDice() int {
	return revealedDiceForSide(b.LastAttack.AttackerRolls, b.LastAttack.DefenderRolls, true, b.RollRevealProgress())
}

// RevealedDefenderDice returns how many defender dice are visible during the roll animation.
func (b *Board) RevealedDefenderDice() int {
	return revealedDiceForSide(b.LastAttack.AttackerRolls, b.LastAttack.DefenderRolls, false, b.RollRevealProgress())
}

func revealedDiceForSide(attacker, defender []int, attackerSide bool, progress float64) int {
	attackerCount := len(attacker)
	defenderCount := len(defender)
	total := attackerCount + defenderCount
	if total == 0 {
		return 0
	}
	if progress >= 1 {
		if attackerSide {
			return attackerCount
		}
		return defenderCount
	}

	revealed := int(progress*float64(total+1) + 0.0001)
	if revealed > total {
		revealed = total
	}
	if attackerSide {
		return min(revealed, attackerCount)
	}

	defenderRevealed := revealed - attackerCount
	if defenderRevealed < 0 {
		return 0
	}
	return min(defenderRevealed, defenderCount)
}
