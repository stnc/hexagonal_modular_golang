package domain

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"gopkg.in/go-playground/validator.v9"
)

func ValidationErrorJson(c fiber.Ctx, err error) error {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_request",
		})
	}

	details := make([]fiber.Map, 0, len(ve))
	for _, fe := range ve {
		details = append(details, fiber.Map{
			"field": fe.Field(),
			"rule":  fe.Tag(),
			"param": fe.Param(),
		})
	}

	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"error":   "validation_failed",
		"details": details,
	})
}
