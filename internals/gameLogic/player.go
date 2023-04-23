package gamelogic

type Player struct {
	ID            uint
	totalScore    uint
	hand          Hand
	isInThisRound bool
}

func NewPlayer(id uint) Player {
	return Player{
		ID:            id,
		totalScore:    0,
		hand:          Hand{},
		isInThisRound: true,
	}
}
