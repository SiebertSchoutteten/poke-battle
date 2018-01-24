package main

import (
	"encoding/csv"
	"fmt"
	"log"
	//"io/ioutil"
	"os"

	"github.com/SiebertSchoutteten/poke-battle/calculator"
)

func main() {
	//log.SetOutput(ioutil.Discard)
	calculator := calculator.NewCalculator()

	//calculator.Fight(poke4, poke2)
	//log.Println(poke3, poke1)

	actualcombats := [][]string{}
	for i := 0; i < 5; i++ {
		poke1 := calculator.GetRandomPokemon()
		poke2 := calculator.GetRandomPokemonWithLevelDifference(poke1.Level(), 10)
		winner := calculator.Fight(poke1, poke2)
		moves1 := poke1.Moves()
		moves2 := poke2.Moves()

		// difference in effectivity of moves
		effectivitydiff := 0.0
		// powerdifference between 2 pokemons moves
		powerdiff := 0
		// ratingdifference between 2 pokemons moves
		ratingdiff := 0.0
		for i := 0; i < len(moves1); i++ {
			move1 := calculator.GetMove(moves1[i])
			move2 := calculator.GetMove(moves2[i])
			if move1 != nil{
				powerdiff += move1.Power
				effectivitydiff += calculator.GetTypeEffectiveness(poke2,move1)
				ratingdiff += move1.Rating
			}
			if move2 != nil{
				powerdiff -= move2.Power
				effectivitydiff -= calculator.GetTypeEffectiveness(poke1,move2)
				ratingdiff -= move2.Rating
			}
		}

		// difference in base stats between the pokemon
		bs1 := poke1.BaseStats()
		bs2 := poke2.BaseStats()
		for i := 0; i < len(bs1); i++ {
			bs1[i] -= bs2[i]
		}

		leveldif := poke1.Level() - poke2.Level()
		poke1wins := false
		if winner == poke1 {
			poke1wins = true
		}
		//if i > 1000000{
			combat := []string{poke1.Name(), poke2.Name(), fmt.Sprintf("%d", leveldif), fmt.Sprintf("%d",bs1[0]), fmt.Sprintf("%d",bs1[1]), fmt.Sprintf("%d",bs1[2]), fmt.Sprintf("%d",bs1[3]),fmt.Sprintf("%d",bs1[4]), fmt.Sprintf("%d",powerdiff),fmt.Sprintf("%d",effectivitydiff),fmt.Sprintf("%d",ratingdiff),fmt.Sprintf("%t", poke1wins)}
			//combat := []string{fmt.Sprintf("%s",poke1.HotEncoding()), fmt.Sprintf("%s",poke2.HotEncoding()), fmt.Sprintf("%d", leveldif),fmt.Sprintf("%t", poke1wins)}
			actualcombats = append(actualcombats, combat)
		//}
	}
	//writeNewCombats(actualcombats)
	//log.Println(poke4.HotEncoding())
	//calculator.GetRandomSpecificPokemon(25,81)
}

// Add Exeception for Bide https://bulbapedia.bulbagarden.net/wiki/Bide_(move)#Effect
// ADD STAB
func writeNewCombats(lines [][]string) {
	file, err := os.Create("combats_test_big_super.csv")
	if err != nil {
		log.Fatalln("error creating new csv")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for _, value := range lines {
		writer.Write(value)
	}
}
