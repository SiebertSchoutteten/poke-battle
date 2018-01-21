package calculator

import "log"
import "fmt"

// Pokemon is short for Pocket Monster
type Pokemon struct {
	base    *PokeBase
	level   int
	moveset [4]*Move
	//ingame status
	status PokemonStatus
	//stages between -6 and 6
	attackBoost      int
	defenseBoost     int
	speedBoost       int
	specialBoost     int
	evasivenessBoost int
	accuracyBoost    int
	//generated stats for base, level and random IV and EV
	stats BaseStats
}

// Attack attacks and returns true if the enemy pokemon dies
func (poke *Pokemon) Attack(enemyMove *Move, enemy *Pokemon, effectiveness float64) bool {
	log.Printf("%s uses %s", enemy.base.Name, enemyMove.Name)
	var attack, defense int
	damage := (2 * enemy.level) / 5
	//log.Println("step 1: ", damage)
	damage += 2
	//log.Println("step 2: ", damage)
	switch enemyMove.Category {
	case "physical":
		attack = enemy.stats.Attack
		defense = poke.stats.Defense
	case "special":
		attack = enemy.stats.Special
		defense = poke.stats.Special
	case "status":
		poke.stats.Hp--
		if poke.stats.Hp <= 0 {
			log.Printf("%s dies", poke.base.Name)
			return true
		}
		return false
	}
	//log.Println("base power:", enemyMove.Power)
	//log.Println("attacck: ", attack)
	//log.Println("defe: ", defense)
	ada := float64(attack) / float64(defense)
	damage *= int(float64(enemyMove.Power) * ada)
	//log.Println("step 3: ", damage)
	damage /= 50
	//log.Println("step 4: ", damage)
	damage += 2
	//log.Println("step 5: ", damage)
	//modifier = random * stab * type effect
	//random is between 0.85 and 1
	modifier := float64(random(218, 255))
	modifier /= 255
	//stab
	for i := 0; i < len(enemy.base.Types); i++ {
		if enemy.base.Types[i] == string(enemyMove.MoveType) {
			modifier *= 1.5
		}
	}
	//type effect

	//modified damage calculation
	damage = int(modifier * float64(damage))
	//log.Println("step 6: ", damage)
	log.Println("actual damage: ", damage)
	poke.stats.Hp -= damage
	if poke.stats.Hp <= 0 {
		log.Printf("%s dies", poke.base.Name)
		return true
	}
	log.Printf("%s has %d hp left", poke.base.Name, poke.stats.Hp)
	return false
}

// Name is a method
func (poke *Pokemon) Name() string {
	return poke.base.Name
}

// Level is a method
func (poke *Pokemon) Level() int {
	return poke.level
}

// Moves is a method
func (poke *Pokemon) Moves() []string {
	return []string{poke.moveset[0].Name, poke.moveset[1].Name, poke.moveset[2].Name, poke.moveset[3].Name}
}

// Stats is a function
func (poke *Pokemon) Stats() []int {
	return []int{poke.stats.Attack, poke.stats.Defense, poke.stats.Hp, poke.stats.Special, poke.stats.Speed}
}

// Number is a methodical function
func (poke *Pokemon) Number() string {
	return fmt.Sprintf("%d", poke.base.Pokenumber)
}

// HotEncoding returns the HotEncoding of a pokemon
func (poke *Pokemon) HotEncoding() string {
	return poke.base.HotEncoding
}
