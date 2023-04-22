package gamelogic

type Player struct {
	ID            uint
	totalScore    uint
	hand          Hand
	isInThisRound bool
}
