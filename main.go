package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"flag"

	"github.com/SiebertSchoutteten/poke-battle/calculator"
)


// PokeMs are
type PokeMs []calculator.PokeM

func main() {
	// Flags
	var dataset bool
	var amount int
	var deck string

	flag.StringVar(&deck, "deck","","provide a pokemon-deck with max 6 pokemon (as in the specified example.json) you want to compare to an opponent to calculate which pokemon has the biggest chance at winning.")
	flag.BoolVar(&dataset, "dataset", false, "provide to generate a file dataset.csv with the amount-param specified amount of simulated battle data")
	flag.IntVar(&amount, "amount", 10000, "when generating a dataset , the amount of data specified will be provided by this param (default: 10000, min: 6)")
	flag.Parse()
	calculato := calculator.NewCalculator()

	if deck != ""{
		// read poke json
		 p := PokeMs{}
		 err := readJSON(deck,&p)
		 if err != nil{
			 log.Fatalln(err)
		 }

		 // check
		 for i := 0; i < len(p); i++ {
			 if !calculato.IsPokemon(p[i].Name){
				 log.Fatalln("error with provided json file")
			 }
			if !calculato.IsMove(p[i].Move1) || !calculato.IsMove(p[i].Move2) || !calculato.IsMove(p[i].Move3) || !calculato.IsMove(p[i].Move4){
				log.Fatalln("error with provided json file")
			}
		 }
		// calculate differences between 6 pkmn and opponent
		allCombats := [][]string{}
		 for i := 1; i < 7; i++ {
			p1 := calculato.GetSpecificPokemon(p[i].Name,p[i].Move1,p[i].Move2,p[i].Move3,p[i].Move4,p[i].Level,p[i].Attack,p[i].Defense,p[i].Hp,p[i].Special,p[i].Speed)
			p2:= calculato.GetSpecificPokemon(p[0].Name,p[0].Move1,p[0].Move2,p[0].Move3,p[0].Move4,p[0].Level,p[0].Attack,p[0].Defense,p[0].Hp,p[0].Special,p[0].Speed)
			combat := calculato.OutputPokemonDifference(p1,p2)
			combat = append(combat, p1.Name())
			allCombats = append(allCombats, combat)
		}
		
		writeNewCombats(allCombats,"combats.csv")
		// provide data to python
		cmd := exec.Command("python", "combat_prediction.py", "combats.csv")
		out, err := cmd.CombinedOutput()
		if err != nil { log.Fatalln(err) }
		// return results 
		log.Println(string(out))
	}

	if dataset{
		if amount > 5 {
			log.SetOutput(ioutil.Discard)			
			actualcombats := [][]string{}
			for i := 0; i < amount; i++ {
				log.Println("fight: ",i)
				poke1 := calculato.GetRandomPokemon()
				poke2 := calculato.GetRandomPokemonWithLevelDifference(poke1.Level(), 99)

				winner := calculato.Fight(poke1,poke2)
				poke1wins := false
				if winner == poke1 {
					poke1wins = true
				}
				//if i > 100000{
					combat := calculato.OutputPokemonDifference(poke1,poke2)
					combat = append(combat, fmt.Sprintf("%t", poke1wins))
					//combat := []string{fmt.Sprintf("%s",poke1.HotEncoding()), fmt.Sprintf("%s",poke2.HotEncoding()), fmt.Sprintf("%d", leveldif),fmt.Sprintf("%t", poke1wins)}
					actualcombats = append(actualcombats, combat)
				//}
			}
			writeNewCombats(actualcombats, "dataset.csv")
			//log.Println(poke4.HotEncoding())
			//calculator.GetRandomSpecificPokemon(25,81)
		}
	}
}

func readJSON(uri string, sort interface{}) error {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		return err
	}
	//m := new(Dispatch)
	//var m interface{}
	err = json.Unmarshal(file, sort)
	if err != nil {
		return err
	}
	return nil
}
func writeNewCombats(lines [][]string, filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660);
	if err != nil {
		log.Fatalln("error creating new csv", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for _, value := range lines {
		writer.Write(value)
	}
}
