package http

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"hexagonalapp/internal/modules/user/app"
	"hexagonalapp/internal/modules/user/domain"
	"time"
)

type APIHandler struct {
	service *app.Service
}

func NewAPI(service *app.Service) *APIHandler {
	return &APIHandler{service: service}
}

func (h *APIHandler) Register(r fiber.Router) {
	r.Post("/users", h.CreateUser)
	r.Get("/users", h.ListUsers)
	r.Get("/users/:id", h.GetUser)
	r.Get("/healthz", h.healthz)
}

func (h *APIHandler) healthz(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "fiber-v3-rest-api",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *APIHandler) CreateUser(c fiber.Ctx) error {
	var req domain.CreateUserInput
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := domain.UsecaseValidate(req); err != nil {
		return domain.ValidationErrorJson(c, err)
	}

	user, validationError, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		return badRequestJson(c, validationError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": user})

}

func (h *APIHandler) GetUser(c fiber.Ctx) error {
	user, err := h.service.GetUser(context.Background(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func (h *APIHandler) ListUsers(c fiber.Ctx) error {
	users, err := h.service.ListUsers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}


func badRequestJson(c fiber.Ctx, message fiber.Map) error {
	return c.Status(fiber.StatusBadRequest).JSON(message)
}
