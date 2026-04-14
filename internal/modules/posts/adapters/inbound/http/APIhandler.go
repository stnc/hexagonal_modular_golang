package http

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"hexagonalapp/internal/modules/posts/app"
)

type APIHandler struct {
	service *app.Service
}

func NewAPI(service *app.Service) *APIHandler {
	return &APIHandler{service: service}
}



func (h *APIHandler) Register(r fiber.Router) {
	r.Post("/posts", h.CreatePost)
	r.Get("/posts", h.ListPosts)
	r.Get("/posts/:id", h.GetPost)
	r.Get("/posts/user/:user_id", h.ListByUser)
}

func (h *APIHandler) CreatePost(c fiber.Ctx) error {
	var req createPostRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.UserID == "" || req.Title == "" || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id, title and content are required"})
	}

	post, err := h.service.CreatePost(context.Background(), app.CreatePostInput{UserID: req.UserID, Title: req.Title, Content: req.Content})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(post)
}

func (h *APIHandler) GetPost(c fiber.Ctx) error {
	post, err := h.service.GetPost(context.Background(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(post)
}

func (h *APIHandler) ListPosts(c fiber.Ctx) error {
	posts, err := h.service.ListPosts(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}

func (h *APIHandler) ListByUser(c fiber.Ctx) error {
	posts, err := h.service.ListPostsByUser(context.Background(), c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}
