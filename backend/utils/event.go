package utils

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/helnr/tubba-game/backend/types"
)

func SendErrorMessage(conn *websocket.Conn, err error, errorType string) error { 
	payloadError, err := json.Marshal(&types.ErrorEventPayload{Error: err.Error(), Type: errorType})
	if err != nil {
		return err
	}
	payload := &types.EventMessage{
		Type: types.ErrorEvent,
		Payload: payloadError,
	}
	conn.WriteJSON(payload)
	return nil
}

func SendNavigateErrorMessage(conn *websocket.Conn, erro error) error { 
	err := SendErrorMessage(conn, erro, "navigate")
	if err != nil {
		return err
	}
	return nil
}