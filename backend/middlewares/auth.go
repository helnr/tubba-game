package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/helnr/tubba-game/backend/types"
	"github.com/helnr/tubba-game/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AuthMiddleware(userStore types.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {

		cookie := c.Cookies("session_id")
		if cookie == "" {
			return utils.WriteUnauthorizedError(c) 
		}

		userID, err := bson.ObjectIDFromHex(cookie)
		if err != nil {
			return utils.WriteUnauthorizedError(c)
		}

		user, err := userStore.GetUserByID(primitive.ObjectID(userID))
		if err != nil {
			return utils.WriteUnauthorizedError(c)
		}
		
		if user.Session.ExpiresAt < time.Now().Unix() || user.Session.IsAborted {
			return utils.WriteUnauthorizedError(c)
		}

		c.Locals("user", user)

		return c.Next()
	}
}