package models

import (
	messages "github.com/CoolCodeTeam/CoolSupportBackend/messages/models"
)

type WebsocketMessage struct {
	WebsocketEventType int              `json:"event_type"`
	Body               messages.Message `json:"body"`
}
