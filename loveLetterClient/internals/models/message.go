package models

const (
	InfoMessage     = "info"
	UpdateMessage   = "update"
	InitDrawMessage = "initDraw"
	TurnDrawMessage = "turnDraw"
	PlayedMessage   = "played"
	AckMessage      = "ack"
)

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type ClientAction struct {
	PlayedCardNumber uint `json:"played_card_number"`
	TargetPlayerId   uint `json:"target_player_id"`
}

type ClientMessage struct {
	ClientId uint
	Message  string
}
