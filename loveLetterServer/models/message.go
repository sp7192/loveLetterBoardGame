package models

type MessageType string

const (
	InfoMessage     MessageType = "info"
	UpdateMessage               = "update"
	InitDrawMessage             = "initDraw"
	TurnDrawMessage             = "turnDraw"
	PlayedMessage               = "played"
	AckMessage                  = "ack"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}
