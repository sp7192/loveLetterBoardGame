package gamelogic

type Player struct {
	ID            uint
	hand          Hand
	totalScore    uint
	isInThisRound bool
}
