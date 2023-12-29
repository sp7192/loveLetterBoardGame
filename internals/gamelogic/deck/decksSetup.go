package deck

import (
	"fmt"
	"io/ioutil"
	"loveLetterBoardGame/internals/gamelogic/card"

	"gopkg.in/yaml.v2"
)

type DeckType string

const (
	Normal   DeckType = "normal"
	Extended DeckType = "extended"
)

type DeckSetup struct {
	Type  string      `yaml:"type"`
	Cards []card.Card `yaml:"cards"`
}

type DeckSetups []DeckSetup

func (ds DeckSetups) FindByType(t DeckType) (DeckSetup, error) {
	for _, v := range ds {
		if v.Type == string(t) {
			return v, nil
		}
	}
	return DeckSetup{}, fmt.Errorf("%s not found", string(t))
}

func ParseDeckSetups(filename string) ([]DeckSetup, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var deckSetups []DeckSetup

	err = yaml.Unmarshal(fileContent, &deckSetups)
	if err != nil {
		return nil, err
	}

	return deckSetups, nil
}
