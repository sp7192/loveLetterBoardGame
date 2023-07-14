package models

const (
	InfoMessage     = "info"
	UpdateMessage   = "update"
	InitDrawMessage = "initDraw"
	TurnDrawMessage = "turnDraw"
	PlayedMessage   = "played"
)

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type ClientAction struct {
	PlayedCardNumber uint `json:"played_card_number"`
}
