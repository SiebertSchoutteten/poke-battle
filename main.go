package main

import (
	"bitbucket.org/be-mobile/pokemon-calculator/calculator"
)


func main() {
	calculator := calculator.NewCalculator()

	poke1 := calculator.GetRandomPokemon()
	poke2 := calculator.GetRandomPokemon()
	calculator.Fight(poke1, poke2)
}

// Add Exeception for Bide https://bulbapedia.bulbagarden.net/wiki/Bide_(move)#Effect
// ADD STAB

