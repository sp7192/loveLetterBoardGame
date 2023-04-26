package gamelogic

import "loveLetterBoardGame/internals/gamelogic/deck"

type GameLogic struct {
	players []Player
	deck    deck.Deck
}

func NewGameLogic(mode string, players []Player) GameLogic {

	return GameLogic{
		players: players,
	}
}
