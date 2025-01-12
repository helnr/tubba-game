package user

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/helnr/tubba-game/backend/config"
	"github.com/helnr/tubba-game/backend/middlewares"
	"github.com/helnr/tubba-game/backend/services/auth"
	"github.com/helnr/tubba-game/backend/types"
	"github.com/helnr/tubba-game/backend/utils"
)

type handler struct {
	userStore types.UserStore

}

func NewHandler(userStore types.UserStore) *handler {
	return &handler{
		userStore: userStore,
	}
}

func (h *handler) RegisterRoutes(fiber *fiber.App) {
	auth := fiber.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)

	user := fiber.Group("/user")
	user.Use(middlewares.AuthMiddleware(h.userStore))
	user.Get("/me", h.GetMe)
	user.Post("/logout", h.Logout)
}

func (h *handler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*types.User)
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data": user,
	})
}

func (h *handler) Register(c *fiber.Ctx) error {
	var payload types.UserRegisterPayload

	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteBadRequestError(c, fmt.Errorf("Invalid payload"))
	}

	err := config.Validator.Struct(&payload)
	if err != nil {
		return utils.WriteValidationError(c, err)
	}

	hashed_password, err := auth.HashPassword(payload.Password)
	if err != nil {
		return utils.WriteInternalServerError(c)
	}



	_, err = h.userStore.GetUserByEmail(payload.Email)
	if err == nil {
		return utils.WriteError(c, fiber.StatusConflict, fmt.Errorf("User with email %s already exists", payload.Email))
	}

	exp := time.Now().Add(time.Hour * 24 * 7).Unix()

	user := types.User{
		Name: payload.Name,
		Email: payload.Email,
		Password: hashed_password,
		Session: types.UserSession{
			ExpiresAt: exp,
			IsAborted: false,
		},
	}

	if err := h.userStore.SaveUser(&user); err != nil {
		return utils.WriteInternalServerError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name: "session_id",
		Value: user.ID.Hex(),
		Expires: time.Unix(exp, 0),
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": user,
	})

}

func (h *handler) Login(c *fiber.Ctx) error {

	cookie := c.Cookies("session_id")
	if cookie != "" {
		return utils.WriteError(c, fiber.StatusUnauthorized, fmt.Errorf("User already logged in"))
	}

	var payload types.UserLoginPayload

	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteBadRequestError(c, fmt.Errorf("Invalid payload"))
	}

	user, err := h.userStore.GetUserByEmail(payload.Email)
	if err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, fmt.Errorf("User not found"))
	}

	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		return utils.WriteError(c, fiber.StatusUnauthorized, fmt.Errorf("Invalid credentials"))
	}


	exp := time.Now().Add(time.Hour * 24 * 7).Unix()

	user.Session.ExpiresAt = exp
	user.Session.IsAborted = false

	if err := h.userStore.UpdateUser(user); err != nil {
		return utils.WriteInternalServerError(c)
	}


	c.Cookie(&fiber.Cookie{
		Name: "session_id",
		Value: user.ID.Hex(),
		Expires: time.Unix(exp, 0),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

func (h *handler) Logout(c *fiber.Ctx) error {

	user := c.Locals("user").(*types.User)

	user.Session.ExpiresAt = 0 
	user.Session.IsAborted = true	

	if err := h.userStore.UpdateUser(user); err != nil {
		return utils.WriteInternalServerError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name: "session_id",
		Value: "",
		Expires: time.Unix(0, 0),
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}