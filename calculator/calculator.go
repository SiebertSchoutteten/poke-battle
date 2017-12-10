package calculator

import (
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"math/rand"
)

type calculator struct {
	pokemon map[int]*PokeBase
	moves   map[int]*Move
	typeEffects map[string]float64
}
// NewCalculator returns a new loaded calculator
func NewCalculator() *calculator {
	calc := &calculator{}
	calc.readPokemon()
	calc.readMoves()
	calc.readEffects()
	return calc
}
func (c *calculator) readPokemon() {

	var pokes Pokebases
	err := readJSON("calculator/pokemon.json", &pokes)
	if err != nil{
		log.Fatalln(err)
	}
	log.Println("read pokemon: ", len(pokes))
	pokemap := make(map[int]*PokeBase)
	for i := 0; i < len(pokes); i++ {
		poke := &pokes[i]
		log.Println(poke)
		pokemap[poke.Number] = poke
	}
	c.pokemon = pokemap
}	
func (c *calculator) readMoves() {
	var moves AllMoves
	err := readJSON("calculator/moves.json", &moves)
	if err != nil{
		log.Fatalln(err)
	}
	movemap := make(map[int]*Move)
	for i := 0; i < len(moves); i++ {
		log.Println(&moves[i])
		movemap[i+1] = &moves[i] 
		movemap[i+1].MoveID = i + 1
 	}
	c.moves = movemap
}
func (c *calculator) readEffects(){
	var typeEffects AllEffectivenesses
	err := readJSON("calculator/type-effectiveness.json", &typeEffects)
	if err != nil{
		log.Fatalln(err)
	}
	typemap := make(map[string]float64)
	for i := 0; i < len(typeEffects); i++ {
		typeformat := fmt.Sprintf("%s-%s", typeEffects[i].Attack, typeEffects[i].Defense)
		typemap[typeformat] = typeEffects[i].Effectiveness
	}
	c.typeEffects = typemap
}
func (c *calculator) GetRandomPokemon() *Pokemon{
	base := c.pokemon[random(0,len(c.pokemon))]
	level := random(MINLEVEL,MAXLEVEL)
	moveset := c.generateMoveset(base, level)
	pokemon := &Pokemon{
		base: base,
		level: level,
		moveset: moveset,
	}

	return pokemon
}
func (c *calculator) generateMoveset(poke *PokeBase, level int) [4]*Move{
	var moves [4]*Move

	for i := 0; i < 4; i++ {
		random := random(1,len(c.moves))
		moves[i] = c.moves[random]
	}
	return moves
}
func (c *calculator) Fight(poke1 *Pokemon, poke2 *Pokemon){
	log.Println("Good Evening ladies & gentlemen")
	fmt.Printf("%+v\n", poke1.base)
	fmt.Printf("%s\n", poke1.level)
	for i := 0; i < 4; i++ {
		fmt.Printf("%+v\n", poke1.moveset[i].Name)
	}
	log.Println("---------------------")
	log.Println("---------------------")
	fmt.Printf("%+v\n", poke2.base)
	fmt.Printf("%s\n", poke2.level)
	for i := 0; i < 4; i++ {
		fmt.Printf("%+v\n", poke2.moveset[i].Name)
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
func random(min, max int) int {
    return rand.Intn(max - min) + min
}