package models

type ClientMessage struct {
	ClientId uint
	Message  string
}

type ClientAction struct {
	PlayedCardNumber uint `json:"played_card_number"`
}
