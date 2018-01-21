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
	for i := 0; i < 1; i++ {
		poke1 := calculator.GetRandomPokemon()
		poke2 := calculator.GetRandomPokemon()
		winner := calculator.Fight(poke1, poke2)
		moves1 := poke1.Moves()
		moves2 := poke2.Moves()
		leveldif := poke1.Level() - poke2.Level()
		poke1wins := false
		if winner == poke1 {
			poke1wins = true
		}
		if i > 1000000{
			combat := []string{poke1.Name(), poke2.Name(), fmt.Sprintf("%d", leveldif), moves1[0], moves1[1], moves1[2], moves1[3], moves2[0], moves2[1], moves2[2], moves2[3], fmt.Sprintf("%t", poke1wins)}
			//combat := []string{fmt.Sprintf("%s",poke1.HotEncoding()), fmt.Sprintf("%s",poke2.HotEncoding()), fmt.Sprintf("%d", leveldif),fmt.Sprintf("%t", poke1wins)}
			actualcombats = append(actualcombats, combat)
		}

	}
	//writeNewCombats(actualcombats)
	//log.Println(poke4.HotEncoding())
	//calculator.GetRandomSpecificPokemon(25,81)
}

// Add Exeception for Bide https://bulbapedia.bulbagarden.net/wiki/Bide_(move)#Effect
// ADD STAB
func writeNewCombats(lines [][]string) {
	file, err := os.Create("combats_test_big.csv")
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
