package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/helnr/tubba-game/backend/config"
)

func WriteError(c *fiber.Ctx, code int, err error) error {
	return c.Status(code).JSON(fiber.Map{
		"status": "error",
		"error": err.Error(),
	})
}

func WriteValidationError(c *fiber.Ctx, err error) error {
	errs := err.(validator.ValidationErrors)

	trans, _ := config.Validator.GetTranslator("en")

	errors := ""
	for _, value := range errs.Translate(trans) {
		errors = fmt.Sprintf("%s, %s", errors, value)
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status": "error",
		"error": errors[2:],
	})
}

func WriteBadRequestError(c *fiber.Ctx, err error) error {
	return WriteError(c, fiber.StatusBadRequest, err)
}

func WriteUnauthorizedError(c *fiber.Ctx) error {
	return WriteError(c, fiber.StatusUnauthorized, fmt.Errorf("Unauthorized"))
}

func WriteInternalServerError(c *fiber.Ctx) error {
	return WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("Something went wrong"))
}