package gamelogic

import "loveLetterBoardGame/cmd/server/internal/gamelogic/card"

type GameState struct {
	DeckCardsCount  uint       `json:"deck_cards_count"`
	PlayersIdInGame []uint     `json:"players_id_in_game"`
	PlayingPlayerId uint       `json:"playing_player_id"`
	PlayedCard      *card.Card `json:"played_card"`
}
