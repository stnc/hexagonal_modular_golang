package http

import (
	"strconv"

	"hexagonalapp/internal/modules/posts/app"
	"hexagonalapp/internal/modules/posts/domain"

	"github.com/gofiber/fiber/v3"
)

type APIHandler struct {
	service *app.Service
}

func NewAPI(service *app.Service) *APIHandler {
	return &APIHandler{service: service}
}

func (h *APIHandler) Register(r fiber.Router) {
	r.Post("/post", h.CreatePost)
	r.Get("/posts", h.ListPosts)
	r.Get("/allposts", h.AllListPosts)
	r.Get("/post/:id", h.GetPost)
	r.Delete("/post/:id", h.DeletePost)
	r.Get("/posts/user/:user_id", h.ListByUser)
}

func (h *APIHandler) CreatePost(c fiber.Ctx) error {
	var req domain.CreatePostInput
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if err := domain.UsecaseValidate(req); err != nil {
		return domain.ValidationErrorJson(c, err)
	}

	post, validationError, err := h.service.CreatePost(c.Context(), req)
	if err != nil {
		return badRequestJson(c, validationError)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": post})
}

func (h *APIHandler) GetPost(c fiber.Ctx) error {
	post, err := h.service.GetPost(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(post)
}

func (h *APIHandler) ListPosts(c fiber.Ctx) error {
	posts, err := h.service.ListPosts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	page := parsePositiveInt(c.Query("page", "1"), 1)
	limit := parsePositiveInt(c.Query("limit", "10"), 10)
	totalItems := len(posts)
	offset := (page - 1) * limit

	if offset >= totalItems {
		return c.JSON(fiber.Map{
			"data": []domain.Post{},
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total_items": totalItems,
			},
		})
	}

	end := offset + limit
	if end > totalItems {
		end = totalItems
	}

	return c.JSON(fiber.Map{
		"data": posts[offset:end],
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_items": totalItems,
		},
	})
}
func (h *APIHandler) AllListPosts(c fiber.Ctx) error {
	posts, err := h.service.ListPosts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(posts)
}

func (h *APIHandler) ListByUser(c fiber.Ctx) error {
	posts, err := h.service.ListPostsByUser(c.Context(), c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}

func (h *APIHandler) DeletePost(c fiber.Ctx) error {
	if _, err := h.service.GetPost(c.Context(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
	}

	if err := h.service.DeletePost(c.Context(), c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func badRequestJson(c fiber.Ctx, message fiber.Map) error {
	return c.Status(fiber.StatusBadRequest).JSON(message)
}

func parsePositiveInt(v string, fallback int) int {
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}
