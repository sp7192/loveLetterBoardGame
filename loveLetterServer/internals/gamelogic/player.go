package gamelogic

import "fmt"

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

func (p *Player) RemoveFromHand(cardId uint) error {
	for i, card := range p.hand.cards {
		if card.Number == cardId {
			p.hand.cards = append(p.hand.cards[:i], p.hand.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card id : %d not in player with id : %d hands", cardId, p.ID)
}
