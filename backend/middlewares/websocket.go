package middlewares

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/helnr/tubba-game/backend/types"
)

func WebsocketMiddleware(games types.GameStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		gameIDString := c.Params("id")
		game, err := games.GetGameByID(gameIDString)
		if err != nil {
			c.Locals("allowed", false)
			c.Locals("error", fmt.Errorf("Game not found"))
			return c.Next()
		}

		c.Locals("gameID", game.ID)
		c.Locals("allowed", true)
		return c.Next()
	}
}