package utils

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/helnr/tubba-game/backend/types"
)


func SendErrorEvent(conn *websocket.Conn, err error) error { 
	payloadError, err := json.Marshal(&types.ErrorEventPayload{Error: err.Error()})
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