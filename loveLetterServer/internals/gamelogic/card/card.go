package card

type Card struct {
	Number uint        `json:"card_number"`
	Effect ICardEffect `json:"-"`
}

func NewCardsSet(mode string) []Card {
	switch mode {
	case "TEST":
		return newTestCardsSet()
	}
	return nil
}

func newTestCardsSet() []Card {
	ret := make([]Card, 0, 100)
	ret = append(ret,
		Card{1, nil},
		Card{1, nil},
		Card{1, nil},
		Card{1, nil},
		Card{2, nil},
		Card{2, nil},
		Card{2, nil},
		Card{3, nil},
		Card{3, nil},
		Card{4, nil},
	)
	return ret
}
