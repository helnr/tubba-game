package types

import "encoding/json"

const (
	GameEvent = "game_event"
	UserEvent = "user_event"
	ErrorEvent = "error_event"
	JoinEvent = "join_event"
	LeaveEvent = "leave_event"
)

type EventMessage struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ErrorEventPayload struct {
	Error string `json:"error"`
}