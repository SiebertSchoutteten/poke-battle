package calculator

// Pokemon is short for Pocket Monster
type Pokemon struct {
	base         *PokeBase
	level        int
	hp           int // pokebase.HP x level
	moveset      [4]*Move
	status       *PokemonStatus
	attackBoost  int
	defenseBoost int
	speedBoost   int
	specialBoost int
}

// Move is an attack that a Pokemon can perform when fighting another Pokemon
type Move struct {
	Name 			 string
	MoveID			  int
	MoveType          PokemonType `json: "type"`
	Category          string
	Power			  int
	PP                int
	Accuracy          float64
	HighCriticalHitRatio bool
	Priority		int
	Effect				interface{}
	Effects				interface{}
	Changes				interface{}
	/*selfDamage        int
		multipleHit       multipleStats
	multipleTurn      multipleStats
	damageRestore     int
	flinchChance      int
	paralyzeChance    int
	poisonChance      int
	sleepChance       int
	burnChance        int
	freezeChance      int
	attackBoostSelf   int
	defenseBoostSelf  int
	specialBoostSelf  int
	speedBoostSelf    int
	attackBoostOther  int
	defenseBoostOther int
	specialBoostOther int
	speedBoostOther   int
	boostChance       int*/
}

// Pokebases are
type Pokebases []PokeBase
// AllMoves are
type AllMoves []Move
// AllEffectivenesses are
type AllEffectivenesses []TypeEffectiveness
// PokeBase defines the base stats for a pokemon, its number, name and 1 or 2 types
// a possible moveset is also part of a PokeBase 
type PokeBase struct {
	Number  int `json: "number"` // Uniquely defined 1-151
	Name    string `json: "name"`
	BaseStats BaseStats `json: "baseStats"`
	Types []string `json: "types"`
	Moveset Moveset `json: "moveset, omitempty"`
}
// BaseStats define a lot
type BaseStats struct{
	Hp      int
	Speed   int
	Special int
	Attack  int
	Defense int
}
// Moveset defines all Moves a Pokemon is able to learn
// Either through leveling or learning from a TM
type Moveset struct {
	movesByLevel []LevelMove
	tmSet []Move
}

// LevelMove specifies a move that a pokemon is able to learn at level
type LevelMove struct {
	level int
	move Move
}

// TypeEffectiveness describes the effectiveness of type attack on type defense
type TypeEffectiveness struct{
	Attack string
	Defense string
	Effectiveness float64
}
type multipleStats struct {
	minMoves int // 2
	hits     []int
}