package game

import "math/rand/v2"

func RollDice(rng *rand.Rand, count int) ([]int, int) {
	rolls := make([]int, count)
	total := 0
	for i := 0; i < count; i++ {
		rolls[i] = rng.IntN(6) + 1
		total += rolls[i]
	}
	return rolls, total
}

type AttackResult struct {
	AttackerRolls []int
	DefenderRolls []int
	AttackTotal   int
	DefenseTotal  int
	AttackerWins  bool
}

func ResolveAttack(rng *rand.Rand, attackerDice, defenderDice int) AttackResult {
	attackerRolls, attackTotal := RollDice(rng, attackerDice)
	defenderRolls, defenseTotal := RollDice(rng, defenderDice)

	return AttackResult{
		AttackerRolls: attackerRolls,
		DefenderRolls: defenderRolls,
		AttackTotal:   attackTotal,
		DefenseTotal:  defenseTotal,
		AttackerWins:  attackTotal > defenseTotal,
	}
}

// attackResultFromRolls rebuilds an AttackResult from recorded rolls, for replay.
func attackResultFromRolls(attackerRolls, defenderRolls []int) AttackResult {
	attackTotal := sumInts(attackerRolls)
	defenseTotal := sumInts(defenderRolls)
	return AttackResult{
		AttackerRolls: attackerRolls,
		DefenderRolls: defenderRolls,
		AttackTotal:   attackTotal,
		DefenseTotal:  defenseTotal,
		AttackerWins:  attackTotal > defenseTotal,
	}
}

func sumInts(vals []int) int {
	total := 0
	for _, v := range vals {
		total += v
	}
	return total
}
