package game

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/helnr/tubba-game/backend/types"
)

type ClientList map[string]map[string]*Client

func (c ClientList) GetGameLength(gameID string) int {
	return len(c[gameID])
}

type Client struct {
	conn *websocket.Conn
	manager *GameHandler
	gameID string
	game *types.Game
	user *types.User
	egress chan types.EventMessage
}

func NewClient(conn *websocket.Conn, manager *GameHandler, gameID string, game *types.Game, user *types.User) *Client {
	return &Client{
		conn: conn,
		manager: manager,
		gameID: gameID,
		game: game,
		user: user,
		egress: make(chan types.EventMessage),
	}
}


func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var eventMessage types.EventMessage
		err = json.Unmarshal(message, &eventMessage)
		if err != nil {
			log.Println(err)
			continue
		}

		err = c.manager.RouteEvent(eventMessage, c)	
		if err != nil {
			log.Println(err)
			continue
		}
		

		// log.Println(len(c.manager.clients[c.gameID]), "clients in game", c.gameID)
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println(err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			err = c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}

func (c *Client) Boadcast(message types.EventMessage) {
	for client_id := range c.manager.clients[c.gameID] {
		client := c.manager.clients[c.gameID][client_id]
		if client != c {
			client.egress <- message
		}
	}
}