package deck

import (
	"loveLetterBoardGame/internals/gamelogic/card"
	"math/rand"
)

type Deck struct {
	cards []card.Card
	count int
}

func NewDeck(cards []card.Card) Deck {
	return Deck{cards: cards, count: len(cards)}
}

func (d *Deck) Draw() (card.Card, bool) {
	if d.count <= 0 {
		return card.Card{}, false
	}
	ret := d.cards[d.count-1]
	d.count--
	return ret, true
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}
