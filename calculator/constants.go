package calculator

// PokemonType is a type that defines a pokemons
type PokemonType int32

// AttackCategory is a category given to a move
type AttackCategory int32

// PokemonStatus is a status a pokemon can reside in
type PokemonStatus int32

// Bug, Dragon, Ice, Fighting, Fire, Flying, Grass, Ghost, Ground, Electric, Normal, Poison, Psychic, Rock, Water
// are all PokemonTypes
const (
	//1, 1.5, 2, 2.5, 3, 3.5, 4, .25, .28, .33, .4, .5, .66
	MAXLEVEL = 99
	MINLEVEL = 1

	Bug      PokemonType = iota
	Dragon   PokemonType = iota
	Ice      PokemonType = iota
	Fighting PokemonType = iota
	Fire     PokemonType = iota
	Flying   PokemonType = iota
	Grass    PokemonType = iota
	Ghost    PokemonType = iota
	Ground   PokemonType = iota
	Electric PokemonType = iota
	Normal   PokemonType = iota
	Poison   PokemonType = iota
	Psychic  PokemonType = iota
	Rock     PokemonType = iota
	Water    PokemonType = iota

	Special  AttackCategory = iota
	Physical AttackCategory = iota
	Status   AttackCategory = iota
	Bide     AttackCategory = iota

	Fit PokemonStatus = iota
	FrozenSolid PokemonStatus = iota
	Burned      PokemonStatus = iota
	Paralyzed   PokemonStatus = iota
	Poisoned    PokemonStatus = iota
	Asleep		PokemonStatus = iota

	Flinched 	PokemonStatus = iota
	Confused	PokemonStatus = iota
	Bound		PokemonStatus = iota
	Leeched		PokemonStatus = iota
)