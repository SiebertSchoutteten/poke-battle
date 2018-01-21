package calculator

//import "log"
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
	return fmt.Sprintf("%d", poke.base.Number)
}
// Speed returns the pokemons speed
func (poke *Pokemon) Speed() int{
	return poke.stats.Speed
}

// HotEncoding returns the HotEncoding of a pokemon
func (poke *Pokemon) HotEncoding() string {
	return poke.base.HotEncoding
}
