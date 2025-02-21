package utils

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/helnr/tubba-game/backend/config"
)

func LogIfDevelopment(msg string) {
	if config.Env.GOEnv == "development" {
		log.Println(msg)
	}
}

func LogErrorIfDevlopment(err error) {
	LogIfDevelopment(err.Error())
}

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

func WriteNotFoundError(c *fiber.Ctx, err error) error {
	return WriteError(c, fiber.StatusNotFound, err)
}

func WriteUnauthorizedError(c *fiber.Ctx) error {
	return WriteError(c, fiber.StatusUnauthorized, fmt.Errorf("Unauthorized"))
}

func WriteInternalServerError(c *fiber.Ctx) error {
	return WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("Something went wrong"))
}