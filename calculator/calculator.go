package calculator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
)

// Calculator knows pokemon
type Calculator struct {
	pokemon     map[int]*PokeBase
	moves       map[int]*Move
	typeEffects map[string]float64
}

// NewCalculator returns a new loaded calculator
func NewCalculator() *Calculator {
	calc := &Calculator{}
	calc.readPokemon()
	calc.readMoves()
	calc.readEffects()
	return calc
}
func (c *Calculator) readPokemon() {

	var pokes Pokebases
	err := readJSON("calculator/pokemon.json", &pokes)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("read pokemon: ", len(pokes))

	pokesList := []string{}
	for i := 0; i < len(pokes); i++ {
		pokesList = append(pokesList, pokes[i].Name)
	}
	hotEncodedList := c.hotEncodeList(pokesList)

	pokemap := make(map[int]*PokeBase)
	for i := 0; i < len(pokes); i++ {
		pokes[i].HotEncoding = hotEncodedList[pokes[i].Name]
		poke := &pokes[i]
		//log.Println(poke.Number)
		pokemap[poke.Number] = poke
	}
	c.pokemon = pokemap
}
func (c *Calculator) readMoves() {
	var moves AllMoves
	err := readJSON("calculator/moves.json", &moves)
	if err != nil {
		log.Fatalln(err)
	}

	movesList := []string{}
	for i := 0; i < len(moves); i++ {
		movesList = append(movesList, moves[i].Name)
	}
	hotEncodedList := c.hotEncodeList(movesList)

	movemap := make(map[int]*Move)
	for i := 0; i < len(moves); i++ {
		//log.Println(&moves[i])
		moves[i].HotEncoding = hotEncodedList[moves[i].Name]
		movemap[i+1] = &moves[i]
		movemap[i+1].MoveID = i + 1
	}
	c.moves = movemap

}
func (c *Calculator) readEffects() {
	var typeEffects AllEffectivenesses
	err := readJSON("calculator/type-effectiveness.json", &typeEffects)
	if err != nil {
		log.Fatalln(err)
	}
	typemap := make(map[string]float64)
	for i := 0; i < len(typeEffects); i++ {
		//log.Println(fmt.Sprintf("%s-%s", typeEffects[i].Attack, typeEffects[i].Defend))
		typeformat := fmt.Sprintf("%s-%s", typeEffects[i].Attack, typeEffects[i].Defend)
		typemap[typeformat] = typeEffects[i].Effectiveness
	}
	c.typeEffects = typemap
}

// GetRandomSpecificPokemon generates a given pokemon with given level
func (c *Calculator) GetRandomSpecificPokemon(pokenumber, level int) *Pokemon {
	log.Println("generating random pokemon")

	//get pokemon base
	base := c.pokemon[pokenumber]
	log.Println("its a ", base.Name)

	//get random level between 1 and 99
	log.Println("with level: ", level)

	//generate a random ABLE moveset for the give pokemon
	moveset := c.generateMoveset(base, level)
	for i := 0; i < len(moveset); i++ {
		log.Println("has move: ", moveset[i].Name)
	}

	//calculate stats for base, level and generated IV's and EV's
	//calculate random EV's
	var attackEV, defenseEV, speedEV, specialEV, HPEV, max int
	if level < 5 {
		max = level
	} else {
		max = ((level - 3) * (level - 3)) + 1
	}
	defeated := random(0, max)
	log.Printf("This pokemon has defeated %d pokemons already", defeated)
	for index := 0; index < defeated; index++ {
		randdefeated := c.pokemon[random(1, 151)]
		attackEV += randdefeated.BaseStats.Attack
		defenseEV += randdefeated.BaseStats.Defense
		speedEV += randdefeated.BaseStats.Speed
		specialEV += randdefeated.BaseStats.Special
		HPEV += randdefeated.BaseStats.Hp
	}

	log.Printf("EV attack: %d, defense: %d, hp: %d, speed: %d, special: %d", attackEV, defenseEV, HPEV, speedEV, specialEV)
	if HPEV > 65535 {
		HPEV = 65535
	}

	if defenseEV > 65535 {
		defenseEV = 65535
	}

	if attackEV > 65535 {
		attackEV = 65535
	}
	if specialEV > 65535 {
		specialEV = 65535
	}

	if speedEV > 65535 {
		speedEV = 65535
	}
	//calculate IV's
	attackIV := random(0, 15)
	defenseIV := random(0, 15)
	speedIV := random(0, 15)
	specialIV := random(0, 15)
	//hpIV is calculated
	var hpIV int
	if !even(attackIV) {
		hpIV += 8
	} // 1111 = 8 + 4 + 2 + 1 = 15
	if !even(defenseIV) {
		hpIV += 4
	}
	if !even(speedIV) {
		hpIV += 2
	}
	if !even(specialIV) {
		hpIV++
	}
	log.Printf("IV attack: %d, defense: %d, hp: %d, speed: %d, special: %d", attackIV, defenseIV, hpIV, speedIV, specialIV)
	stats := &BaseStats{
		Hp:      c.calculateHP(base.BaseStats.Hp, hpIV, HPEV, level),
		Attack:  c.calculateOtherStat(base.BaseStats.Attack, attackIV, attackEV, level),
		Defense: c.calculateOtherStat(base.BaseStats.Defense, defenseIV, defenseEV, level),
		Speed:   c.calculateOtherStat(base.BaseStats.Speed, speedIV, speedEV, level),
		Special: c.calculateOtherStat(base.BaseStats.Special, specialIV, specialEV, level),
	}

	log.Printf("attack: %d, defense: %d, hp: %d, speed: %d, special: %d", stats.Attack, stats.Defense, stats.Hp, stats.Speed, stats.Special)

	pokemon := &Pokemon{
		base:    base,
		level:   level,
		moveset: moveset,
		stats:   *stats,
		status:  Fit,
		recharging: false,
		volatileStatus: make(map[PokemonStatus]bool),
		maxHP: stats.Hp,
	}

	return pokemon
}

// GetRandomPokemon generates a random pokemon with a random level and a random possible moveset
func (c *Calculator) GetRandomPokemon() *Pokemon {
	log.Println("generating random pokemon")

	//get pokemon base
	base := c.pokemon[random(1, 151)]
	log.Println("its a ", base.Name)

	//get random level between 1 and 99
	level := random(MINLEVEL, MAXLEVEL)
	log.Println("with level: ", level)

	//generate a random ABLE moveset for the give pokemon
	moveset := c.generateMoveset(base, level)
	for i := 0; i < len(moveset); i++ {
		log.Println("has move: ", moveset[i].Name)
	}

	//calculate stats for base, level and generated IV's and EV's
	//calculate random EV's
	var attackEV, defenseEV, speedEV, specialEV, HPEV, max int
	if level < 5 {
		max = level
	} else {
		max = ((level - 3) * (level - 3)) + 1
	}
	max /= 3
	defeated := random(0, max)
	log.Printf("This pokemon has defeated %d pokemons already", defeated)
	for index := 0; index < defeated; index++ {
		randdefeated := c.pokemon[random(1, 151)]
		attackEV += randdefeated.BaseStats.Attack
		defenseEV += randdefeated.BaseStats.Defense
		speedEV += randdefeated.BaseStats.Speed
		specialEV += randdefeated.BaseStats.Special
		HPEV += randdefeated.BaseStats.Hp
	}
	if HPEV > 65535 {
		HPEV = 65535
	}

	if defenseEV > 65535 {
		defenseEV = 65535
	}

	if attackEV > 65535 {
		attackEV = 65535
	}
	if specialEV > 65535 {
		specialEV = 65535
	}

	if speedEV > 65535 {
		speedEV = 65535
	}
	//calculate IV's
	attackIV := random(0, 15)
	defenseIV := random(0, 15)
	speedIV := random(0, 15)
	specialIV := random(0, 15)
	//hpIV is calculated
	var hpIV int
	if !even(attackIV) {
		hpIV += 8
	} // 1111 = 8 + 4 + 2 + 1 = 15
	if !even(defenseIV) {
		hpIV += 4
	}
	if !even(speedIV) {
		hpIV += 2
	}
	if !even(specialIV) {
		hpIV++
	}

	stats := &BaseStats{
		Hp:      c.calculateHP(base.BaseStats.Hp, hpIV, HPEV, level),
		Attack:  c.calculateOtherStat(base.BaseStats.Attack, attackIV, attackEV, level),
		Defense: c.calculateOtherStat(base.BaseStats.Defense, defenseIV, defenseEV, level),
		Speed:   c.calculateOtherStat(base.BaseStats.Speed, speedIV, speedEV, level),
		Special: c.calculateOtherStat(base.BaseStats.Special, specialIV, specialEV, level),
	}

	log.Printf("attack: %d, defense: %d, hp: %d, speed: %d, special: %d", stats.Attack, stats.Defense, stats.Hp, stats.Speed, stats.Special)

	pokemon := &Pokemon{
		base:    base,
		level:   level,
		moveset: moveset,
		stats:   *stats,
		status:  Fit,
		recharging: false,
		volatileStatus: make(map[PokemonStatus]bool),
		maxHP: stats.Hp,
	}

	return pokemon
}
func (c *Calculator) generateMoveset(poke *PokeBase, level int) [4]*Move {
	var moves [4]*Move

	for i := 0; i < 4; i++ {
		random := random(1, len(c.moves))
		moves[i] = c.moves[random]
	}
	return moves
}

// Returns whether poke1 attacks first or not
func (c *Calculator) poke1First(poke1Speed, poke2Speed, poke1Priority, poke2Priority int) bool{

		// The pokemon that selected the move with the highest priority will attack first
		if poke1Priority > poke2Priority {
			return true
		}
		if poke2Priority > poke1Priority{
			return false
		}
		//  if both moves have the same priority, the pokemon with the higher speed will attack first
		if poke1Speed > poke2Speed {
			return true
		}
		if poke2Speed > poke1Speed {
			return false
		}
		// if both moves have the same priority and both pokemon have the same speed
		// it is determined randomly who will attack first
		if random(0, 1) == 0 {
			return true
		}

		return false
}

// Fight simulates a fight between poke1 and poke2
func (c *Calculator) Fight(poke1 *Pokemon, poke2 *Pokemon) *Pokemon {

	// As long as one of the pokemon hasnt reached 0 HP the fight isnt over yet
	for {
		//let poke1 and poke2 choose a random move before fighting, unless a move was used that will be continued to use 
		poke1move := poke1.SelectMove()
		poke2move := poke2.SelectMove()
		

		//if the chosen move is metronome a random move will be chosen
		for poke1move.Name == "metronome"{
			poke1move = c.moves[random(0,len(c.moves)-1)]
		}
		for poke2move.Name == "metronome"{
			poke2move = c.moves[random(0,len(c.moves)-1)]
		}

		if poke1move.Name == "mirror move"{
			if poke2.lastMove != nil{
				poke1move = poke2.lastMove
			}
		}
		if poke2move.Name == "mirror move"{
			if poke1.lastMove != nil{
				poke2move = poke1.lastMove
			}
		}

		//it is then decided in this method which pokemon will attack first, paralyzed pokemon operate at 25% speed
		poke1speed := poke1.Speed()
		poke2speed := poke2.Speed()

		if poke1.status == Paralyzed{
			poke1speed /= 4
		}

		if poke2.status == Paralyzed{
			poke2speed /= 4
		}
		
	
		poke1first := c.poke1First(poke1speed, poke2speed,poke1move.Priority, poke2move.Priority)
		

		// In case pokemon 2 uses counter
		var is1physical,is2physical bool
		if poke1move.MoveType == "normal" ||poke1move.MoveType == "fighting"{
			is1physical = true
		}
		if poke2move.MoveType == "normal" ||poke2move.MoveType == "fighting"{
			is2physical = true
		}
		// attacking happens here
		// if an returns a pokemon it means that pokemon died (could be target or attacker)
		// if a pokemon misses its attack, nothing happens
		// if a pokemon is unable to attack due to paralyzation or other status effects nothing happens either
		if poke1first {
			someonedied := c.TryToAttack(poke2, poke1, poke1move,is1physical,0)
			poke1.lastMove = poke1move
			if someonedied != nil{
				return someonedied
			}
			someonedied = c.TryToAttack(poke1, poke2, poke2move,is2physical,poke2.lastDealtDamage)
			poke2.lastMove = poke2move
			if someonedied != nil{
				return someonedied
			}
		} else {
			someonedied := c.TryToAttack(poke1, poke2, poke2move,is2physical,0)
			poke2.lastMove = poke2move
			if someonedied != nil{
				return someonedied
			}
			someonedied = c.TryToAttack(poke2, poke1, poke1move,is1physical,poke1.lastDealtDamage)
			poke1.lastMove = poke1move
			if someonedied != nil{
				return someonedied
			}
		}

		// After effects applysion happens here, one of both pokemon could die after applysion
		didsomeonedie := c.ApplyAfterEffects(poke1, poke2)
		if didsomeonedie != nil{
			return didsomeonedie
		}
	}
}

// TryToAttack lets an attacker attack the target if possible, it returns the target/attacker if one of them died, nil if both survived
func (c *Calculator) TryToAttack(target, attacker *Pokemon, usedmove *Move, physical bool,damage int) *Pokemon{
	if c.isPokemonAbleToAttack(attacker) && target.invulnerability == false{
		if attacker.volatileStatus[Confused] && (random(1,100) < 50){
			log.Println("The pokemon is confused, it hurt itself in its confusion")
			attacker.stats.Hp -= 8
			if attacker.stats.Hp <= 0{
				return attacker
			}
		}else{
			if c.doesMoveHit(usedmove,attacker.accuracyBoost,target.evasivenessBoost){
				// Pokemon 1's attack strikes, if attack returns true, pokemon1 wins
				effectiveness := c.getTypeEffectiveness(target, usedmove)
				dmg :=  c.Attack(usedmove, target, attacker, effectiveness,physical,damage)
				if target.IsDead(){
					return target
				}

				if attacker.IsDead(){
					return attacker
				}
				c.Effect(usedmove, dmg, target, attacker)
				if target.IsDead(){
					return target
				}

				if attacker.IsDead(){
					return attacker
				}
			}else{
				log.Println("Attack Missed")
			}
		}
	}
	return nil
}
// ApplyAfterEffects applies after battle effects such as burn or poison, and returns a Pokemon if one died, nil if none died
func (c *Calculator) ApplyAfterEffects(poke1, poke2 *Pokemon) *Pokemon{
			// status effects such as poison and burn happen here
			if poke1.status == Burned ||  poke1.status == Poisoned{
				poke1.stats.Hp -= (poke1.maxHP / 16)
				log.Println("burned or poisoned")
			}
			if poke2.status == Burned ||  poke2.status == Poisoned{
				poke2.stats.Hp -= (poke2.maxHP / 16)
				log.Println("burned or poisoned")
			}

			// sleep wears off after 1-4 turns
			if poke1.status == Asleep{
				if random(1,100) < 50 || poke1.sleepTurn == 4{
					log.Println("the pokemon woke up")
				}else{
					poke1.sleepTurn++
				}
			}

			if poke2.status == Asleep{
				if random(1,100) < 50 || poke2.sleepTurn == 4{
					log.Println("the pokemon woke up")
				}else{
					poke2.sleepTurn++
				}
			}

			// Non volatile status effects
			if poke1.volatileStatus[Bound]{
				if random(1,100) < 34 || poke1.boundTurn == 4 {
					poke1.Unbind()
					log.Println("The pokemon was released")
				}else{
					poke1.boundTurn++
					poke1.stats.Hp -= (poke1.maxHP / 16)
					log.Println("poke was hurt by bind")
				}
			}

			if poke2.volatileStatus[Bound]{
				if random(1,100) < 34 || poke2.boundTurn == 4 {
					poke2.Unbind()
					log.Println("The pokemon was released")
				}else{
					poke2.boundTurn++
					poke2.stats.Hp -= (poke2.maxHP / 16)
					log.Println("poke was hurt by bind")
				}
			}

			// confusion might wear off
			if poke1.volatileStatus[Confused]{
				if poke1.confusedTurn == 4 || random(1,100) < 50{
					log.Println("The pokemon snapped out of its confusion")
					poke1.Unconfuse()
				}
			}

			if poke2.volatileStatus[Confused]{
				if poke2.confusedTurn == 4 || random(1,100) < 50{
					log.Println("The pokemon snapped out of its confusion")
					poke2.Unconfuse()
				}
			}

			// a leeched pokemon has 1/16th of its max HP drained to the opponent
			if poke1.volatileStatus[Leeched]{
				damage := poke1.maxHP / 16
				poke1.stats.Hp -= damage
				poke2.Lifesteal(damage)
			}

			if poke2.volatileStatus[Leeched]{
				damage := poke2.maxHP / 16
				poke2.stats.Hp -= damage
				poke1.Lifesteal(damage)
			}


			// flinch is only active until the end of a turn
			poke1.Unflinch()
			poke2.Unflinch()

			// A disabled move might wear off after max 6 turns
			if poke1.disabledMoveTurn == 6 || random(0,100) < 41{
				log.Println("not disabled anymore")
				poke1.disabledMove = nil
				poke1.disabledMoveTurn = 0
			}

			if poke2.disabledMoveTurn == 6 || random(0,100) < 41{
				log.Println("not disabled anymore")
				poke2.disabledMove = nil
				poke2.disabledMoveTurn = 0
			}


			// If one of the pokemon is dead return it, else return nil
			if poke1.IsDead(){
				return poke1
			}

			if poke2.IsDead(){
				return poke2
			}

			return nil
}
// RandomCriticalMove returns true if the provided move is a random critical
func (c *Calculator) RandomCriticalMove(poke *Pokemon, move *Move) bool{
	probability := poke.base.BaseStats.Speed / 2

	if poke.critical{
		probability /= 4
	}

	if move.HighCriticalHitRatio{
		probability *= 8
	}

	if random(1,256) < probability{
		return true
	}

	return false
}
// Effect handles the effect that comes with an attack causing one of the following effects:
func (c *Calculator) Effect(move *Move, damage int,poke, enemy *Pokemon){

	randomFactor := random(1,100)
	//effects that only work if substitute is active
	if poke.substituteHp <= 0 {
		// acid has a 10 procent chance of lowering  the targets defense stats with 1
		if move.Name == "acid" && !poke.misted{
			if randomFactor > 11 {
				poke.ModifyDefense(-1)
			}
		}

		// aurora beam has a 10 procent chance of lowering the targets attack stats with 1
		if move.Name == "aurora beam" && !poke.misted {
			if randomFactor < 11 {
				poke.ModifyAttack(-1)
			}
		}
		// bite has a 30 procent chance of flinching the target
		if move.Name == "bite" || move.Name == "headbutt" || move.Name == "rockslide" || move.Name == "rolling kick" || move.Name == "stomp" {
			if randomFactor < 31 {
				poke.Flinch()
			}
		}
		// blizzard has a 10 procent chance of freezing the target
		if move.Name == "blizzard" || move.Name == "ice beam" || move.Name == "ice punch"{
			if randomFactor < 11 {
				poke.Freeze()
			}
		}
		// body slam has a 30 procent chance of paralyzing the target
		if move.Name == "body slam" || move.Name == "lick" || move.Name == "low kick"{
			if randomFactor < 31{ 
				poke.Paralyze()
			}
		}

		// bone club, hyper fang has a 10 procent chance of flinching the target
		if move.Name == "bone club" || move.Name == "hyper fang"{
			if randomFactor < 11{
				poke.Flinch()
			}
		}
		// constrict, bubble, bubblebeam have a 10 procent chance of lowering the targets speed stats with 1
		if (move.Name == "bubble" || move.Name == "bubblebeam" || move.Name == "constrict") && !poke.misted {
			if randomFactor < 11{
				poke.ModifySpeed(-1)
			}
		}
			// ember has a 10 procent chance of burnin the target
		if move.Name == "ember" || move.Name == "fire blast" || move.Name == "fire punch" || move.Name == "flamethrower"{
			if randomFactor > 11 {
				poke.Burn()
			}
		}

		// confuse ray confuses the target
		if move.Name == "confuse ray" || move.Name == "supersonic"{
			poke.Confuse()	
		}
		// flash drops the targets accuracy 1 stage
		if (move.Name == "flash" || move.Name == "kinesis" || move.Name == "sand attack" || move.Name == "smokescreen") && !poke.misted{
			poke.ModifyAccuracy(-1)
		}
		// growl lowers the targets attack stats with 1
		if move.Name == "growl" && !poke.misted{
			poke.ModifyAttack(-1)
		}
		//leer modifies target defense stat with -1
		if (move.Name == "Leer" || move.Name == "tail whip") && !poke.misted{
			poke.ModifyDefense(-1)
		}
		//poison gas poisons the target
		if move.Name == "poison gas" || move.Name == "poison powder" || move.Name == "toxic"{
			poke.Poison()
		}
		// poison sting has 30% chance of poisoning
		if move.Name == "poison sting" || move.Name == "sludge"{
			if randomFactor < 31{
				poke.Poison()
			}
		}
		// screech lowers defense stats with 2
		if move.Name == "screech" && !poke.misted{
			poke.ModifyDefense(-2)
		}
	
		// smog has a 40% chance of poisoning
		if move.Name == "smog"{
			if randomFactor < 41{
				poke.Poison()
			}
		}
		// string shot reduces speed stat with 1
		if move.Name == "string shot" && !poke.misted{
			poke.ModifySpeed(-1)
		}
		if move.Name == "thunder" || move.Name == "thunderbolt" || move.Name == "thunder shock" || move.Name == "thunder punch"{
			if randomFactor < 11 {
				poke.Paralyze()
			}
		}
	}

	// transform transforms the attacking pokemon into the target
	if move.Name == "transform"{
		oldhp := enemy.stats.Hp
		enemy.stats = poke.stats
		enemy.stats.Hp = oldhp

		enemy.moveset = poke.moveset
		enemy.attackBoost = poke.attackBoost
		enemy.defenseBoost = poke.defenseBoost
		enemy.speedBoost = poke.speedBoost
		enemy.specialBoost = poke.specialBoost
		enemy.evasivenessBoost = poke.evasivenessBoost
		enemy.accuracyBoost = poke.accuracyBoost

		enemy.base.Types = poke.base.Types
	}
	// disable disables a random target move
	if move.Name == "disabled" {
		if poke.disabledMove == nil{
			poke.disabledMove = poke.moveset[random(0,3)]
		}else{
			poke.disabledMoveTurn++
		}
	}
	// mist mists some misty mist
	if move.Name == "mist"{
		enemy.misted = true
	}
	// rest heals to 100% HP and sets user to sleep for 2 turns
	if move.Name == "rest"{
		if enemy.recurrentMove.Name == "rest"{
			enemy.recurrentMoveTurn++
			if enemy.recurrentMoveTurn > 1{
				enemy.Cure()
			}
		}else{
			if enemy.stats.Hp != enemy.maxHP{
				if enemy.recurrentMoveTurn == 0{
					enemy.stats.Hp = enemy.maxHP
					enemy.Cure()
					enemy.Sleep()
					enemy.recurrentMove = move
				}	
			}
		}
		
	}

	// substitute creates substituteHP from own HP
	if move.Name == "substitute"{
		if enemy.stats.Hp >= (enemy.maxHP / 4){
			enemy.stats.Hp -= (enemy.maxHP /4)
			enemy.substituteHp = enemy.maxHP + 1
		}
	}

	// transform transforms  
	// reflect enables a reflection
	if move.Name == "reflect"{
		enemy.reflected = true
	}
	// Light screen enables a light screen
	if move.Name == "light screen" {
		enemy.lightscreen = true
	}

	// mimic copies a move of the other user
	if move.Name == "mimic"{
		var index int
		for i := 0; i < 4; i++ {
			if enemy.moveset[i].Name == "mimic"{
				index = i
			}
		}
		enemy.ChangeMove(index, poke.moveset[random(0,3)])
	}

	// struggle deals recoil damage
	if move.Name == "struggle"{
		enemy.stats.Hp -= (damage/2)
	}
	// absorb and dream eater drain 50% from the target to the enemy
	if move.Name == "absorb" || move.Name == "dream eater" || move.Name == "leech life" || move.Name == "mega drain"{
		enemy.Lifesteal(damage/2)
	}



	// acid armor raises the defense stat of a pokemon with 2 stages
	// barrier raises the defense stat of a pokemon with 2
	if move.Name == "acid armor" || move.Name == "barrier" {
		enemy.ModifyDefense(2)
	}

	// agility raises the speed stat of a pokemon with 2 stages
	if move.Name == "agility" {
		enemy.ModifySpeed(2)
	}

	// amnesia raises the special stat of a pokemon with 2 stages
	if move.Name == "amnesia" {
		enemy.ModifySpecial(2)
	}



	// bind, fire spin, clamp bind a pokemon for a number of turns
	if move.Name == "bind" || move.Name == "fire spin" || move.Name == "clamp"{
		poke.Bind()
	}

	


	// confusion has a 10 procent chance of confusing the target
	if move.Name == "confusion" || move.Name == "psybeam"{	
		if randomFactor < 11 {
			poke.Confuse()
		}
	}
	// conversion changes the enemys type to the targets type
	if move.Name == "conversion"{
		enemy.ChangeTypes(poke.base.Types)
	}
	// defense curl, harden raises the pokemons defense stats with 1
	if move.Name == "defense curl" || move.Name == "harden" || move.Name == "withdraw"{
		enemy.ModifyDefense(1)
	}
	// double team raises pokemons evasiveness
	if move.Name == "double team" || move.Name == "minimize"{
		enemy.ModifyEvasiveness(1)
	}
	//double-edge recoils 1/4th hp of the actual damage
	if move.Name == "double-edge" || move.Name == "take down"{
		enemy.stats.Hp -= (damage / 4)
	}

	// glare paralyzes the target
	if move.Name == "glare" || move.Name == "stun spore" || move.Name == "thunder wave"{	
		poke.Paralyze()
	}
		
	// Growth raises the pokemons special stats with 1
	if move.Name == "growth"{
		enemy.ModifySpecial(1)
	}
	// Haze does a lot
	if move.Name == "haze"{
		poke.ResetStats()
		enemy.ResetStats()
		enemy.Cure()

		poke.Unleech()
		enemy.Unleech()

		poke.lightscreen = false
		enemy.lightscreen = false

		poke.reflected = false
		enemy.lightscreen = false

		poke.misted = false
		poke.misted = false

		poke.Unconfuse()
		enemy.Unconfuse()
	}
	// if high jump kick missed, enemy recoils 1 HP
	if move.Name == "high jump kick" || move.Name == "jump kick"{
		if damage == 0{
			enemy.stats.Hp-- 
		}
	}

	// hypnosis lets a pokemon sleep
	if move.Name == "hypnosis" || move.Name == "lovely kiss" || move.Name == "sing" || move.Name == "sleep powder" || move.Name == "spore"{
		poke.Sleep()
	}


	// medidate raises pokemon attack stat with 1
	if move.Name == "meditate" || move.Name == "sharpen"{
		enemy.ModifyAttack(1)
	}
	
	// recover recovers up to 50% of max HP
	if move.Name == "recover" || move.Name == "soft-boiled"{
		hptorecover := random(1,enemy.maxHP/2)
		enemy.stats.Hp += hptorecover

		if enemy.stats.Hp > enemy.maxHP{
			enemy.stats.Hp = enemy.maxHP
		}
	}

	if move.Name == "swords dance"{
		enemy.ModifyAttack(2)
	}


	// Recurrent moves
	if move.Name == "thrash" || move.Name == "skull bash" || move.Name == "solar beam" || move.Name == "hyper beam" || move.Name == "sky attack" || move.Name == "petal dance" || move.Name == "dig" || move.Name == "fly" || move.Name == "bide" || move.Name == "rage" || move.Name == "razor wind"{
		enemy.recurrentMove = move
		enemy.recurrentMoveTurn++
	}

	// bide is reset if it dealt damage
	if move.Name == "bide"{
		enemy.recurrentMove = nil
		enemy.recurrentMoveTurn = 0
		enemy.bideCount = 0
	}

	// rage never stops 
	if move.Name == "rage"{
		if damage > 0 {
			enemy.ModifyAttack(1)
		}
	}

	// thrash and petal dance might wear off after 3-4 turns and then confuse the pokemon
	if move.Name == "thrash" || move.Name == "petal dance"{
		if enemy.recurrentMoveTurn == random(3,4) || enemy.recurrentMoveTurn == 4{
			enemy.recurrentMove = nil
			enemy.recurrentMoveTurn = 0
			enemy.Confuse()
		}
	}
	// these attacks become unconcurrent after 2 turns
	if  move.Name == "skull bash" || move.Name == "solar beam" || move.Name == "hyper beam" || move.Name == "sky attack" || move.Name == "razor wind"{
		if enemy.recurrentMoveTurn == 2{
			enemy.recurrentMove = nil
			enemy.recurrentMoveTurn = 0
		}
	}
	// dig and fly make a pokemon invulnerable for one turn and become unconcurrent after 2 turns
	if move.Name == "dig" || move.Name == "fly"{
		if enemy.recurrentMoveTurn == 1{
			enemy.Invulnerate()
		}
		if enemy.recurrentMoveTurn == 2{
			enemy.recurrentMove = nil
			enemy.recurrentMoveTurn = 0
			enemy.Uninvulnerate()
		}
	}

	// focus energy lowers critical hit chances
	if move.Name == "focus energy"{
		enemy.critical = true
	}

	// leech seed leeches a pokemon
	if move.Name == "leech seed"{
		poke.Leech()
	}
	

	// if the target is biding, add damage to its bidecount
	if poke.recurrentMove.Name == "bide"{
		poke.bideCount += damage
	}
}
// Attack attacks and returns the effective damage the enemy did on the target
func (c *Calculator) Attack(enemyMove *Move, poke, enemy *Pokemon, effectiveness float64, physical bool, phDmg int) int {

	// One hit KO's in case speed is faster
	if enemyMove.Name == "fissure" || enemyMove.Name == "horn drill" || enemyMove.Name == "guillotine" {
		if enemy.stats.Speed > poke.stats.Speed{
			poke.stats.Hp = 0
			return poke.maxHP
		}
	}

	// Selfdestruct and exploision kill the attacker
	if enemyMove.Name == "explosion" || enemyMove.Name == "selfdestruct"{
		enemy.stats.Hp = 0
		return enemy.maxHP
	}

	// decide if move is critical
	critical := c.RandomCriticalMove(enemy, enemyMove)

	// Damage calculation 
	log.Printf("%s uses %s", enemy.base.Name, enemyMove.Name)
	var attack, defense int

	level := enemy.level
	if critical{
		log.Println("critical hit")
		level *= 2
	}

	damage := (2 * level) / 5

	
	log.Println("step 1: ", damage)
	damage += 2
	log.Println("step 2: ", damage)

	switch enemyMove.Category {
	case "physical":
		attack = enemy.stats.Attack
		//burned pokemon's attack stat is halved
		if enemy.status == Burned  && !critical{
			attack /= 2
		}
		defense = poke.stats.Defense
		if !critical{
			attack =  int(float64(attack) *c.getStageMultiplier(false, enemy.attackBoost))
			defense = int(float64(defense) *c.getStageMultiplier(false, poke.defenseBoost))
		}
	case "special":
		attack = enemy.stats.Special
		defense = poke.stats.Special
		if !critical{
			attack =  int(float64(attack) *c.getStageMultiplier(false, enemy.attackBoost))
			if enemy.lightscreen{
				attack *= 2
			}
			defense = int(float64(defense) *c.getStageMultiplier(false, poke.defenseBoost))
			if poke.reflected{
				defense *= 2
			}
		}
	case "status":
		attack = 0
		defense = 0
	}

	ada := int(float64(attack) / float64(defense))
	log.Println("attack/defense", ada)

	log.Println("Power: ", enemyMove.Power)
	damage *= enemyMove.Power * int(ada)
	log.Println("step 3: ", damage)
	damage /= 50
	log.Println("step 4: ", damage)
	damage += 2
	log.Println("step 5: ", damage)
	//modifier = random * stab * type effect

	//random is between 0.85 and 1
	Mrandom := float64(random(218, 255))
	Mrandom /= 255

	//stab
	Mstab := 1.0
	for i := 0; i < len(enemy.base.Types); i++ {
		if enemy.base.Types[i] == string(enemyMove.MoveType) {
			Mstab = 1.5
		}
	}
	log.Println("Stab: ", Mstab)
	// type effect was given in param
	log.Println("effectivness: ", effectiveness)
	// So modifier is calculated
	modifier := Mrandom * Mstab *effectiveness
	log.Println("modifier: ", modifier)

	//modified damage calculation
	damage = int(modifier * float64(damage))
	//log.Println("step 6: ", damage)
	log.Println("actual damage: ", damage)


	turns := 1
	// special turn moves
	if enemyMove.Name == "bonemerang" || enemyMove.Name == "double kick" || enemyMove.Name == "twineedle"{
		turns = 2
	}

	if enemyMove.Name == "barrage" || enemyMove.Name == "comet punch" || enemyMove.Name == "double slap" || enemyMove.Name == "fire spin" || enemyMove.Name == "fury swipes" || enemyMove.Name == "fury attack" || enemyMove.Name == "pin missile" || enemyMove.Name == "spike cannon"{
		turns = random(2,5)
	}
	// Special damage moves
	if enemyMove.Name == "psywave"{
		damage = random(int(enemy.Level()/2), int(enemy.Level()*(3/2)))
	}
	if enemyMove.Name == "seismic toss" || enemyMove.Name == "night shade"{
		damage = enemy.Level()
	}
	if enemyMove.Name == "dragon rage"{
		damage = 40
	}

	if enemyMove.Name == "super fang"{
		damage = int(poke.stats.Hp / 2) - 1
	}

	if enemyMove.Name == "sonic boom"{
		damage = 20
	}

	if enemyMove.Name == "dream eater" && poke.status != Asleep{
		damage = 0
	}
	// Recurrent moves that dont deal damage in case its not the right turn
	if enemyMove.Name == "solar beam" || enemyMove.Name == "sky attack" || enemyMove.Name == "skull bash" || enemyMove.Name == "fly" || enemyMove.Name == "dig" || enemyMove.Name == "razor wind"{
		if enemy.recurrentMoveTurn > 1{
			log.Println("charging")
			damage = 0
		}
	}

	if enemyMove.Name == "hyper beam"{
		if enemy.recurrentMoveTurn == 1{
			log.Println("recharging")
			damage = 0
		}
	}

	if enemyMove.Name == "fly" || enemyMove.Name == "dig"{
		if enemy.recurrentMoveTurn > 1{
			log.Println("underground or high in the sky")
			damage = 0
		}
	}

	if enemyMove.Name == "bide"{
		if enemy.recurrentMoveTurn == 3 || enemy.recurrentMoveTurn == random(2,3){
			damage = enemy.bideCount * 2
		}else{
			damage = 0
		}
	}

	if enemyMove.Name == "counter"{
		if physical{
			damage = phDmg * 2
		}else{
			damage = 0
		}	
	}

	// here pokemon deals damage x turns
	log.Println("Applied damage:", damage)
	for i := 0; i < turns; i++ {
		if i > 1{
			log.Println("turn :", i)
		}
		if poke.substituteHp > 0 && enemyMove.Name != "super fang" && enemyMove.Name != "bind" && enemyMove.Name != "clamp" && enemyMove.Name != "fire spin"{
			poke.substituteHp -= damage
			if poke.substituteHp < 0 {
				// if substitute breaks all multi turns moves break
				log.Println("Substitute broke")
				poke.substituteHp = 0
				i = 6
			}
		}else{
			poke.stats.Hp -= damage
		}
	}

	log.Printf("%s has %d hp left", poke.base.Name, poke.stats.Hp)
	return damage
}
func (c *Calculator) getTypeEffectiveness(poke *Pokemon, move *Move) float64 {
	effectiveness := c.typeEffects[string(move.MoveType)+"-"+poke.base.Types[0]]

	if len(poke.base.Types) > 1 {
		effectiveness *= c.typeEffects[string(move.MoveType)+"-"+poke.base.Types[1]]
	}

	return effectiveness
}
func (c *Calculator) isPokemonAbleToAttack(poke *Pokemon) bool {

	if poke.status == Asleep || poke.status == FrozenSolid || poke.status == Flinched{
		log.Println("Flinched, asleep or frozen solid")
		return false
	}
	if poke.status == Paralyzed{
		if random(1,100) > 75 {
			log.Println("Paralyzed")
			return false
		}
	}
	return true
}
func (c *Calculator) doesMoveHit(move *Move, selfAccuracyStage, enemyEvasionStage int) bool {
	//always hits
	if move.Name == "swift" {
		return true
	}
	// 4, 1
	selfMultiplier := c.getStageMultiplier(true, selfAccuracyStage)
	enemyMultiplier := c.getStageMultiplier(false, enemyEvasionStage)
	
	//1/256 glitch if multipliers are unchanged 
	if selfMultiplier == 1 && enemyMultiplier == 1 {
		if random(1,256) == 256{
			log.Println("glitch")
			return false
		}
	}

	hitratio := (selfMultiplier * enemyMultiplier) * move.Accuracy
	if hitratio >= 1 {
		return true
	}

	chances := random(1, 100)
	if (hitratio*100)-float64(chances) > 0 {
		return true
	}
	return false
}
func (c *Calculator) getStageMultiplier(accuracy bool, stage int) float64 {
	stageMultipliers := []float64{0.25, 0.28, 0.33, 0.4, 0.5, 0.66, 1, 1.5, 2, 2.5, 3, 3.5, 4}
	// if its not accuracy, its evasion
	if !accuracy {
		stage -= 6
		// stage should now be between -12 and 0
		newStage := math.Abs(float64(stage))
		return stageMultipliers[int(newStage)]
	}
	// stage should be between 0 and 12
	return stageMultipliers[stage+6]
}
func (c *Calculator) calculateHP(baseHP, hpIV, hpEV, level int) int {
	//stats are rounded down if decimal
	sqrtEV := math.Sqrt(float64(hpEV)) / 4
	baseIV := float64(baseHP + hpIV)
	hp := (((baseIV*2 + sqrtEV) * float64(level)) / 100) + float64(level) + 10
	return int(hp)
}
func (c *Calculator) calculateOtherStat(base, iv, ev, level int) int {
	//stats are rounded down if decimal
	baseIV := float64(base + iv)
	sqrtEV := math.Sqrt(float64(ev)) / 4
	//this calculation is rounded down for any stat except HP
	stat := int(baseIV*2 + sqrtEV)
	stat = (stat * level) / 100
	return stat + 5
}

func (c *Calculator) hotEncodeList(list []string) map[string]string {
	encodedList := make(map[string]string, len(list))

	for j := 0; j < len(list); j++ {
		element := "("
		for i := 0; i < len(list); i++ {
			if i != j {
				element += "0"
			} else {
				element += "1"
			}

			if i != (len(list) - 1) {
				element += ","
			}
		}
		element += ")"
		encodedList[list[j]] = element
	}
	return encodedList
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
	if (max - min) <= 0 {
		return 1
	}
	return rand.Intn(max-min) + min
}
func even(number int) bool {
	return number%2 == 0
}
