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

func ResolveAttack(rng *rand.Rand, attackerDice, defenderDice int, cheatAttacker, cheatDefender bool) AttackResult {
	attackCount := effectiveDiceCount(attackerDice, cheatAttacker)
	defenseCount := effectiveDiceCount(defenderDice, cheatDefender)

	attackerRolls, attackTotal := RollDice(rng, attackCount)
	defenderRolls, defenseTotal := RollDice(rng, defenseCount)

	return AttackResult{
		AttackerRolls: attackerRolls,
		DefenderRolls: defenderRolls,
		AttackTotal:   attackTotal,
		DefenseTotal:  defenseTotal,
		AttackerWins:  attackTotal > defenseTotal,
	}
}

func effectiveDiceCount(base int, cheat bool) int {
	count := base
	if cheat {
		count = 1 + int(float64(count)*1.5)
	}
	if count < 1 {
		return 1
	}
	return count
}
