package calculator

import "log"
import "fmt"

// Pokemon is short for Pocket Monster
type Pokemon struct {
	base    *PokeBase
	level   int
	moveset [4]*Move
	pp	 [4]int
	disabledMove	*Move
	disabledMoveTurn	int
	lastMove *Move
	//ingame status
	substituteHp int
	misted bool
	lightscreen bool
	reflected bool
	invulnerability bool
	status PokemonStatus
	volatileStatus map[PokemonStatus]bool
	recurrentMove		*Move
	recurrentMoveTurn	 int
	bideCount		int
	sleepTurn		int
	boundTurn		int
	confusedTurn	int
	//stages between -6 and 6
	attackBoost      int
	defenseBoost     int
	speedBoost       int
	specialBoost     int
	evasivenessBoost int
	accuracyBoost    int
	//generated stats for base, level and random IV and EV
	critical bool
	lastDealtDamage int
	maxHP int
	stats BaseStats
}

func (poke *Pokemon) moveCount() int{
	count := 0
	for i := 0; i < len(poke.moveset); i++ {
		if poke.moveset[i] != nil{
			if poke.moveset[i].Name != "none"{
				count++
			}
		}
	}
	return count
}
// returns all moves that are not out of pp
func (poke *Pokemon) ableMoves() []string{
	amountOfMoves := poke.moveCount()
	moves := []string{}
	log.Println("amount",amountOfMoves)
	for i := 0; i < amountOfMoves; i++ {
		log.Println(poke.moveset[i].Name, poke.pp[i])
		if poke.pp[i] >= 0{
			if poke.disabledMove != nil{
				if poke.disabledMove.Name != poke.moveset[i].Name{
					moves = append(moves, poke.moveset[i].Name)
				}
			}else{
				moves = append(moves, poke.moveset[i].Name)
			}

		}
	}
	return moves
}
// SelectMove decides a random move the pokemon will use (recurrent moves will automatically be selected), disabled moves cannot be selected
func (poke *Pokemon) SelectMove() string{
	if poke.recurrentMove != nil{
		return poke.recurrentMove.Name
	}

	moves := poke.ableMoves()
	log.Println("able moves: ", moves)
	if len(moves) == 0{
		//out of moves
		return "struggle"
	}

	if len(moves) == 1{
		for i := 0; i < len(poke.moveset); i++ {
			if poke.moveset[i].Name == moves[0]{
				poke.pp[i]--
			}
		}
		return moves[0]
	}	
	move := moves[random(0,len(moves)-1)]
	for i := 0; i < len(poke.moveset); i++ {
		if poke.moveset[i].Name == move{
			poke.pp[i]--
		}
	}
	return move
}
// ChangeMove a pokemon's move on index to the given move
func (poke *Pokemon) ChangeMove(index int, move *Move){
	poke.moveset[index] = move
}
// ChangeTypes changes the pokemons type
func (poke *Pokemon) ChangeTypes(types []string){
	poke.base.Types = types
}

// ResetStats resets the pokemons stats including critical chance
func (poke *Pokemon) ResetStats(){
	poke.attackBoost = 0
	poke.defenseBoost = 0
	poke.specialBoost = 0
	poke.speedBoost = 0
	poke.accuracyBoost = 0
	poke.evasivenessBoost = 0
	poke.critical = false
}

// Sleep sleeps the pokemon in if possible
func (poke *Pokemon) Sleep(){
	if poke.status == Fit{
		poke.status = Asleep
	}
}
// Burn burns the pokemon if possible
func (poke *Pokemon) Burn(){
	if poke.status == Fit{
		poke.status = Burned
	}
}

// Bind binds the pokemon
func (poke *Pokemon) Bind(){
	poke.volatileStatus[Bound] = true
}

// Invulnerate invulnerates the pokemon
func (poke *Pokemon) Invulnerate(){
	poke.invulnerability = true
}

// Uninvulnerate invulnerates the pokemon
func (poke *Pokemon) Uninvulnerate(){
	poke.invulnerability = false
}

// Unbind binds the pokemon
func (poke *Pokemon) Unbind(){
	poke.boundTurn = 0
	poke.volatileStatus[Bound] = false
}

// Poison the pokemon if possible
func (poke *Pokemon) Poison(){
	if poke.status == Fit{
		poke.status = Poisoned
	}
}

// Flinch flinches a pokemon
func (poke *Pokemon) Flinch(){
	poke.volatileStatus[Flinched] = true
}

// Unflinch unflinches a pokemon, which happens after every turn
func (poke *Pokemon) Unflinch(){
	poke.volatileStatus[Flinched] = false
}
// Leech leeches the pokemon 
func (poke *Pokemon) Leech(){
	poke.volatileStatus[Leeched] = true
}

// Unleech unleeches a pokemon
func (poke *Pokemon) Unleech(){
	poke.volatileStatus[Leeched] = false
}
// Confuse confuses the pokemon 
func (poke *Pokemon) Confuse(){
	poke.volatileStatus[Confused] = true
}

// Unconfuse unconfuses the pokemon
func (poke *Pokemon) Unconfuse(){
	poke.confusedTurn = 0
	poke.volatileStatus[Confused] = false
}
// Paralyze paralyzes a pokemon if it doesnt have another status condition
func (poke *Pokemon) Paralyze(){
	if poke.status == Fit {
		poke.status = Paralyzed
	}
}

// Freeze freezs a pokemon if it doesnt have another status condition
func (poke *Pokemon) Freeze(){
	if poke.status == Fit {
		poke.status = FrozenSolid
	}
}
// Cure cures a pokemon from any status condition
func (poke *Pokemon) Cure(){
	poke.sleepTurn = 0
	poke.status = Fit
}
// ModifyAttack modifies the specified pokemon with x stages, stages are always between -3 en 3
func (poke *Pokemon) ModifyAttack(stages int){
	if (poke.attackBoost + stages) >= 6 {
		poke.attackBoost = 6
		log.Println("Max stages reached")
	} else if (poke.attackBoost + stages) <= -6{
		poke.attackBoost = -6
		log.Println("Min stages reached")
	}else{
		poke.attackBoost += stages
	}
}
// ModifyDefense modifies the specified pokemon with x stages, stages are always between -3 en 3
func (poke *Pokemon) ModifyDefense(stages int){
	if (poke.defenseBoost + stages) >= 6 {
		poke.defenseBoost = 6
		log.Println("Max stages reached")
	} else if (poke.defenseBoost + stages) <= -6{
		poke.defenseBoost = -6
		log.Println("Min stages reached")
	}else{
		poke.defenseBoost += stages
	}
}
// ModifySpeed modifies the specified pokemon with x stages, stages are always between -3 en 3
func (poke *Pokemon) ModifySpeed(stages int){
	if (poke.speedBoost + stages) >= 6 {
		poke.speedBoost = 6
		log.Println("Max stages reached")
	} else if (poke.speedBoost + stages) <= -6{
		poke.speedBoost = -6
		log.Println("Min stages reached")
	}else{
		poke.speedBoost += stages
	}
}
// ModifySpecial modifies the specified pokemon with x stages, stages are always between -3 en 3
func (poke *Pokemon) ModifySpecial(stages int){
	if (poke.specialBoost + stages) >= 6 {
		poke.specialBoost = 6
		log.Println("Max stages reached")
	} else if (poke.specialBoost + stages) <= -6{
		poke.specialBoost = -6
		log.Println("Min stages reached")
	}else{
		poke.specialBoost += stages
	}
}
// ModifyEvasiveness modifies the specified pokemon with x stages, stages are always between -3 en 3
func (poke *Pokemon) ModifyEvasiveness(stages int){
	if (poke.evasivenessBoost + stages) >= 6 {
		poke.evasivenessBoost = 6
		log.Println("Max stages reached")
	} else if (poke.evasivenessBoost + stages) <= -6{
		poke.evasivenessBoost = -6
		log.Println("Min stages reached")
	}else{
		poke.evasivenessBoost += stages
	}
}
// ModifyAccuracy modifies the specified pokemon with x stages, stages are always between -3 en 3
func (poke *Pokemon) ModifyAccuracy(stages int){
	if (poke.accuracyBoost + stages) >= 6 {
		poke.accuracyBoost = 6
		log.Println("Max stages reached")
	} else if (poke.accuracyBoost + stages) <= -6{
		poke.accuracyBoost = -6
		log.Println("Min stages reached")
	}else{
		poke.accuracyBoost += stages
	}
}
// Lifesteal steals a given amount of life from another pokemon
func (poke *Pokemon) Lifesteal(hp int)  {
	poke.stats.Hp += hp
}
// Types returns a pokemons types
func (poke *Pokemon) Types() []string{
	return poke.base.Types
}
// IsDead returns whether a pokemon is dead or not
func (poke *Pokemon) IsDead() bool{
	if poke.stats.Hp <= 0{
		return true
	}
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

// BaseStats is a function
func (poke *Pokemon) BaseStats () []int{
	return []int{poke.base.BaseStats.Attack,poke.base.BaseStats.Defense,poke.base.BaseStats.Hp,poke.base.BaseStats.Special,poke.stats.Speed}
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
