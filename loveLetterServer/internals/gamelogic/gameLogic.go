package gamelogic

import (
	"encoding/json"
	"fmt"
	"loveLetterBoardGame/internals/gamelogic/card"
	"loveLetterBoardGame/internals/gamelogic/deck"
	"loveLetterBoardGame/models"
	"math/rand"
)

type GameLogic struct {
	Players            []Player
	Deck               deck.Deck
	PlayingPlayerIndex uint
	PlayingPlayerId    uint
}

func NewGameLogic(mode string, players []Player) GameLogic {
	cards := card.NewCardsSet("TEST")
	return GameLogic{
		Deck:    deck.NewDeck(cards),
		Players: players,
	}
}

func (g *GameLogic) getPlayerById(playerId uint) (*Player, error) {
	for _, player := range g.Players {
		if player.ID == playerId {
			return &player, nil
		}
	}
	return nil, fmt.Errorf("Player with %d id not found", playerId)
}

func (g *GameLogic) getStartingPlayerIndex() uint {
	return uint(rand.Intn(len(g.Players)))
}

func (g *GameLogic) PreparePhase() error {
	g.Deck.Shuffle()
	for i := range g.Players {
		card, ok := g.Deck.Draw()
		if !ok {
			return fmt.Errorf("not enough cards in deck")
		}
		g.Players[i].hand.cards = append(g.Players[i].hand.cards, card)
	}
	g.PlayingPlayerIndex = g.getStartingPlayerIndex()
	return nil
}

func (g *GameLogic) isPlayerIdValid(playerId uint) bool {
	for _, p := range g.Players {
		if playerId == p.ID {
			return true
		}
	}
	return false
}

func (g *GameLogic) GetGameState() (GameState, error) {
	var ret GameState

	ret.DeckCardsCount = uint(g.Deck.Count())
	for _, p := range g.Players {
		if p.isInThisRound {
			ret.PlayersIdInGame = append(ret.PlayersIdInGame, p.ID)
		}
	}

	ret.PlayingPlayerId = g.Players[g.PlayingPlayerIndex].ID

	return ret, nil
}

func (g *GameLogic) GetPlayersCardsInHand(id uint) []card.Card {
	for _, p := range g.Players {
		if id == p.ID {
			return p.hand.cards
		}
	}
	return nil
}

func (g *GameLogic) DrawPhase() bool {
	card, ok := g.Deck.Draw()
	if !ok {
		return false
	}
	index := g.PlayingPlayerIndex
	g.Players[index].hand.cards = append(g.Players[index].hand.cards, card)
	return true
}

func (g *GameLogic) UpdateGame(msg models.ClientMessage) error {
	var action models.ClientAction
	err := json.Unmarshal([]byte(msg.Message), &action)
	if err != nil {
		return err
	}

	playingPlayer, err := g.getPlayerById(msg.ClientId)
	if err != nil {
		return err
	}

	playedCard, err := playingPlayer.RemoveFromHand(action.PlayedCardId)
	if err != nil {
		return err
	}

	// TODO : At the moment cards have no effect
	if playedCard.Effect != nil {
		playedCard.Effect.Play()
	}

	return nil
}
