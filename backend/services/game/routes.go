package game

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/helnr/tubba-game/backend/middlewares"
	"github.com/helnr/tubba-game/backend/types"
	"github.com/helnr/tubba-game/backend/utils"
)

var (
	upgrader = gwebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)


type GameHandler struct {
	userStore types.UserStore 
	games types.GameStore
	handlers map[string]EventHandler
	sync.RWMutex
	sync.WaitGroup
}

func NewGameHandler(userStore types.UserStore, gameStore types.GameStore) *GameHandler {
	gameHandler := &GameHandler{
		userStore: userStore,
		games : gameStore,
		handlers: make(map[string]EventHandler),
	}

	gameHandler.setupEventHandlers()
	return gameHandler
}

func (h *GameHandler) setupEventHandlers() {
	h.handlers[types.GameEvent] = HandleGameEvent  
	h.handlers[types.JoinEvent] = HandleJoinEvent
	h.handlers[types.LeaveEvent] = HandleLeaveEvent
	h.handlers[types.ReadyEvent] = HandleReadyEvent
	h.handlers[types.ChangeTeamEvent] = HandleChangeTeamEvent
	h.handlers[types.StartGameEvent] = HandleStartEvent
	h.handlers[types.CardPlayedEvent] = HandleCardPlayedEvent 
	h.handlers[types.TubbaEvent] = HandleTubbasEvent
}

func (h *GameHandler) RegisterRoutes(app fiber.Router) {
	app.Post("/game", middlewares.AuthMiddleware(h.userStore), h.CreateGame)
	app.Get("/game/:id", middlewares.AuthMiddleware(h.userStore), h.GetGame)

	app.Get("/game/join/:id", middlewares.AuthMiddleware(h.userStore), middlewares.WebsocketMiddleware(h.games), websocket.New(h.JoinGame))
}

func (h *GameHandler) RouteEvent(event types.EventMessage, player *types.Player) error {
	if handler, ok := h.handlers[event.Type]; ok {
		return handler(&event, player)	
	}

	return fmt.Errorf("Unknown event")
}
func (h *GameHandler) GetGame(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*types.User)
	

	game, err := h.games.GetGameByID(id)
	if err != nil {
		return utils.WriteBadRequestError(c, fmt.Errorf("Game not found"))
	}

	response := fiber.Map{
		"status": "success",
		"game": game.Info(),
	}


	if game.Status == types.STATUS_LOBBY {
		
		if game.Players.Has(user.ID) {
			return utils.WriteBadRequestError(c, fmt.Errorf("You are already in the game"))
		}
		if game.Players.Len() >= 4 {
			return utils.WriteBadRequestError(c, fmt.Errorf("Game is full"))
		}

	} else {

		if !game.Players.Has(user.ID) {
			return utils.WriteBadRequestError(c, fmt.Errorf("You are not part of the game"))
		}
	}

	

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *GameHandler) CreateGame(c *fiber.Ctx) error {
	game := types.NewGame()	
	gameID, err := h.games.AddGame(game)

	if err != nil {
		return utils.WriteInternalServerError(c)
	}

	game.ID = gameID
	game.OwnerID = c.Locals("user").(*types.User).ID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"game": game.Info(),
	})
}


func (h *GameHandler) JoinGame(conn *websocket.Conn) {
	allowed := conn.Locals("allowed").(bool)

	if !allowed {
		err := conn.Locals("error").(error)
		err = utils.SendNavigateErrorMessage(conn, err)
		if err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}

	user := conn.Locals("user").(*types.User)

	gameIDString := conn.Locals("gameID").(string)
	game, err := h.games.GetGameByID(gameIDString)
	if err != nil {
		err := utils.SendNavigateErrorMessage(conn, fmt.Errorf("Game is full"))
		if err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}


	var player *types.Player

	if game.Status == types.STATUS_LOBBY {

		if game.Players.Has(user.ID) {
			err := utils.SendNavigateErrorMessage(conn, fmt.Errorf("you are already in the game"))
			if err != nil {
				log.Println(err)
			}
			conn.Close()
			return
		}

		if game.Players.Len() >= 4 {
			err := utils.SendNavigateErrorMessage(conn, fmt.Errorf("Game is full"))
			if err != nil {
				log.Println(err)
			}
			conn.Close()
			return
		}

		player = types.NewPlayer(conn, h, gameIDString, game, user)

	}else {
		if !game.Players.Has(user.ID) {
			err := utils.SendNavigateErrorMessage(conn, fmt.Errorf("Game is started"))
			if err != nil {
				log.Println(err)
			}
			conn.Close()
			return
		}

		oldPlayer := game.Players.Get(user.ID)

		if oldPlayer.Connected {
			err := utils.SendNavigateErrorMessage(conn, fmt.Errorf("You are already in the game"))
			if err != nil {
				log.Println(err)
			}
			conn.Close()
			return
		}

		player = types.NewPlayerFromOldPlayer(oldPlayer, conn, h, gameIDString, game, user)
		h.RemovePlayer(oldPlayer)
	}


	if player.User.ID == game.OwnerID {
		player.IsOwner = true
	}

	h.AddPlayer(player)

	// Read messages
	go player.ReadMessages()
	player.WriteMessages()
}

func (h *GameHandler) AddPlayer(player *types.Player) {
	h.Lock()
	defer h.Unlock()

	h.games.AddPlayer(player.GameID, player)
}

func (h *GameHandler) RemovePlayer(player *types.Player) {
	h.Lock()
	defer h.Unlock()

	h.games.RemovePlayer(player.GameID, player)
}

func (h *GameHandler) GetGameByID(id string) (*types.Game, error) {
	return h.games.GetGameByID(id)
}