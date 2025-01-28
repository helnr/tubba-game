package game

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/helnr/tubba-game/backend/middlewares"
	"github.com/helnr/tubba-game/backend/types"
	"github.com/helnr/tubba-game/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	upgrader = gwebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)


type GameHandler struct {
	userStore types.UserStore 
	gameStore types.GameStore
	clients ClientList
	handlers map[string]EventHandler
	sync.RWMutex
	sync.WaitGroup
}

func NewGameHandler(userStore types.UserStore, gameStore types.GameStore) *GameHandler {
	gameHandler := &GameHandler{
		userStore: userStore,
		gameStore: gameStore,
		clients: make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	gameHandler.setupEventHandlers()
	return gameHandler
}

func (h *GameHandler) setupEventHandlers() {
	h.handlers[types.GameEvent] = HandleGameEvent  
	h.handlers[types.JoinEvent] = HandleJoinEvent
	h.handlers[types.LeaveEvent] = HandleLeaveEvent
}

func (h *GameHandler) RegisterRoutes(fiber *fiber.App) {
	fiber.Post("/game", middlewares.AuthMiddleware(h.userStore), h.CreateGame)
	fiber.Get("/game/:id", middlewares.AuthMiddleware(h.userStore), h.GetGame)

	fiber.Get("/game/join/:id", middlewares.AuthMiddleware(h.userStore), middlewares.WebsocketMiddleware(h.gameStore), websocket.New(h.JoinGame))
	fiber.Get("game/test/:id/:name", middlewares.WebsocketMiddleware(h.gameStore), websocket.New(h.TestFiber))
}

func (h *GameHandler) RouteEvent(event types.EventMessage, client *Client) error {
	if handler, ok := h.handlers[event.Type]; ok {
		return handler(&event, client)	
	}

	return fmt.Errorf("Unknown event")
}
func (h *GameHandler) GetGame(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*types.User)


	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return utils.WriteBadRequestError(c, fmt.Errorf("Invalid game ID"))
	}

	game, err := h.gameStore.GetGameByID(primitive.ObjectID(objectID))
	if err != nil {
		return utils.WriteBadRequestError(c, fmt.Errorf("Game not found"))
	}

	var found bool
	for _, player := range game.Players {
		if player.ID == user.ID {
			found = true
			break
		}
	}

	if !found {
		if !(len(game.Players) < 4) {
			return utils.WriteBadRequestError(c, fmt.Errorf("Game is full"))
		}

		game.Players = append(game.Players, types.Player{
			ID: user.ID,
			Cards: []string{},
			Name: user.Name,
			Team: "",
			IsTurn: false,
		})

		if err := h.gameStore.UpdateGame(game); err != nil {
			return utils.WriteInternalServerError(c)
		}
	}

	


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"game": game.Info(),
	})
}

func (h *GameHandler) CreateGame(c *fiber.Ctx) error {
	game := &types.Game{
		Status: types.STATUS_LOBBY,
		Players: []types.Player{},
		CreatedAt: time.Now(),
	}

	err := h.gameStore.SaveGame(game)
	if err != nil {
		return utils.WriteInternalServerError(c)
	}

	gameInfo := &types.GameInfo{
		ID: game.ID,
		Status: game.Status,
		Players: []types.PlayerInfo{},
	}
	

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"game": gameInfo,
	})
}


func (h *GameHandler) JoinGame(conn *websocket.Conn) {
	allowed := conn.Locals("allowed").(bool)

	if !allowed {
		err := conn.Locals("error").(error)
		err = utils.SendErrorEvent(conn, err)
		if err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}

	gameIDString := conn.Params("id")
	if h.clients.GetGameLength(gameIDString) >= 4 {
		err := utils.SendErrorEvent(conn, fmt.Errorf("Game is full"))
		if err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}


	user := conn.Locals("user").(*types.User)
	game := conn.Locals("game").(*types.Game)

	if game.Status != types.STATUS_LOBBY {
		err := utils.SendErrorEvent(conn, fmt.Errorf("Game is started"))
		if err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}

	client := NewClient(conn, h, gameIDString, game, user)

	h.addClient(client)

	// Read messages
	go client.readMessages()
	client.writeMessages()
}

func (h *GameHandler) addClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.clients[client.gameID]; !ok {
		h.clients[client.gameID] = make(map[string]*Client)
	}
	h.clients[client.gameID][client.user.Name] = client
}

func (h *GameHandler) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.clients[client.gameID]; ok {
		client.conn.Close()
		delete(h.clients[client.gameID], client.user.Name)

	}
}

func containes(values []string, value string) bool {
	for _, item := range values {
		if value == item {
			return true
		}
	}
	return false
}

var users = [4]string{"mohammed", "ali", "ahmed", "hasan"}


func (h *GameHandler) Test(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := map[string]string{
			"status": "error",
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(err)
		return
	}

	context := context.WithValue(r.Context(), "conn", conn)
	r = r.WithContext(context)

	conn.Close()
}

func (h *GameHandler) TestFiber(conn *websocket.Conn) {

	user_id := conn.Params("id")
	username := conn.Params("name")
	game := conn.Locals("game").(*types.Game)

	if !containes(users[:], username) {
		conn.Close()
		return
	}

	var user types.User = types.User{
		Name: username,
	}

	client := NewClient(conn, h, user_id, game, &user)

	h.addClient(client)

	go client.readMessages()
	client.writeMessages()

}