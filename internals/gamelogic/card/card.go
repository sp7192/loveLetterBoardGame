package card

type Card struct {
	Number uint        `json:"number" yaml:"number"`
	Effect ICardEffect `json:"-"`
	Name   string      `yaml:"name"`
	Help   string      `yaml:"help"`
	Count  int         `yaml:"count"`
	Type   string      `yaml:"type"`
}

func NewCardsSet(mode string) []Card {
	switch mode {
	case "TEST":
		return testDeckSetup()
	}
	return nil
}

func testDeckSetup() []Card {
	ret := make([]Card, 0, 100)
	ret = append(ret,
		Card{Number: 1},
		Card{Number: 2},
		Card{Number: 3},
		Card{Number: 4},
		Card{Number: 1},
		Card{Number: 2},
		Card{Number: 3},
		Card{Number: 4},
		Card{Number: 5},
		Card{Number: 6},
	)
	return ret
}
