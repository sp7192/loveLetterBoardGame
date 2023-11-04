package models

type MessageType string

const (
	InfoMessage     MessageType = "info"
	UpdateMessage   MessageType = "update"
	InitDrawMessage MessageType = "initDraw"
	TurnDrawMessage MessageType = "turnDraw"
	PlayedMessage   MessageType = "played"
	AckMessage      MessageType = "ack"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}

func (msg MessageType) String() string {
	return string(msg)
}
