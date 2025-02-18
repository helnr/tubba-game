package types

import "encoding/json"

const (
	GameEvent = "game_event"
	UserEvent = "user_event"
	ErrorEvent = "error_event"
	JoinEvent = "join_event"
	LeaveEvent = "leave_event"
	ReadyEvent = "ready_event"
	ChangeTeamEvent = "change_team_event"
	StartGameEvent = "start_event"
	CardPlayedEvent = "played_card_event"
	TubbaEvent = "tubba_event"
)

type EventMessage struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewEventMessage(eventType string, payload any) *EventMessage {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return &EventMessage{
		Type: eventType,
		Payload: payloadBytes,
	}
}

type ErrorEventPayload struct {
	Error string `json:"error"`
	Type string `json:"type"`
}