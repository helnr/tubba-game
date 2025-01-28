package middlewares

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/helnr/tubba-game/backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func WebsocketMiddleware(gameStore types.GameStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		gameIDString := c.Params("id")
		gameID, err := bson.ObjectIDFromHex(gameIDString)
		if err != nil {
			c.Locals("allowed", false)
			c.Locals("error", fmt.Errorf("Invalid game ID"))
			return c.Next()
		}

		game, err := gameStore.GetGameByID(primitive.ObjectID(gameID))
		if err != nil {
			c.Locals("allowed", false)
			c.Locals("error", fmt.Errorf("game no longer exists"))
			return c.Next()
		}

		c.Locals("gameID", gameID)
		c.Locals("game", game)
		c.Locals("allowed", true)
		return c.Next()
	}
}