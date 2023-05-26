package models

type ClientMessage struct {
	ClientId uint
	Message  string
}

type ClientAction struct {
	PlayedCardId uint `json:"played_card_id"`
}
