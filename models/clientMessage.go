package models

type ClientMessage struct {
	ClientId uint
	Message  string
}

type ClientAction struct {
	PlayedCardNumber uint `json:"played_card_number"`
	TargetPlayerId   uint `json:"target_player_id"`
}
