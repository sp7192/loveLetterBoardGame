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

func CreatePlayersFromIDs(ids []uint) []Player {
	players := make([]Player, 0, len(ids))

	for _, id := range ids {
		player := Player{
			ID:            id,
			totalScore:    0,
			hand:          Hand{},
			isInThisRound: true,
		}

		players = append(players, player)
	}

	return players
}
