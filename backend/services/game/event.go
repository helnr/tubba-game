package game

import (
	"encoding/json"
	"log"

	"github.com/helnr/tubba-game/backend/types"
)






type GameEventPayload struct {
	Status types.GameStatus `json:"status"`
}

type JoinEventPayload struct {
	UserID string `json:"user_id"`
	Username string `json:"username"`
}

type EventHandler func(message *types.EventMessage, client *Client) error

func HandleGameEvent(message *types.EventMessage, client *Client) error {
	var payload GameEventPayload
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		return err
	}

	log.Println("Game event: ", payload.Status)
	client.Boadcast(*message)
	return nil
}

func HandleLeaveEvent(message *types.EventMessage, client *Client) error {
	client.manager.removeClient(client)
	log.Println("Leave ", client.manager.clients[client.gameID])
	return nil
}

func HandleJoinEvent(message *types.EventMessage, client *Client) error {
	
	return nil
}

func SendGameEvent(client *Client) error {
	payload := GameEventPayload{
		Status: client.game.Status,
	}
	raw_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	client.Boadcast(types.EventMessage{
		Type: types.GameEvent,
		Payload: raw_payload,
	})

	return nil
}

func SendJoinEvent(client *Client) error {
	payload := JoinEventPayload{
		UserID: client.user.ID.Hex(),
		Username: client.user.Name,
	}
	raw_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	client.Boadcast(types.EventMessage{
		Type: types.JoinEvent,
		Payload: raw_payload,
	})

	return nil
}

