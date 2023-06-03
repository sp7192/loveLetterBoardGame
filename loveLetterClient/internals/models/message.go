package models

const (
	InfoMessage   = "info"
	UpdateMessage = "update"
	DrawMessage   = "draw"
	PlayedMessage = "played"
)

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}
