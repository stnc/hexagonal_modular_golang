package domain

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"hexagonalapp/internal/platform/adapters/inbound/http/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gopkg.in/go-playground/validator.v9"
)

type Post struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostInput struct {
	ID      uint   `json:"id"`
	UserID  string `json:"user_id"`
	Title   string `json:"title" validate:"required,min=2,max=200"`
	Content string `json:"content" validate:"required,min=2,max=5000"`
}

func NormalizeInput(input CreatePostInput) CreatePostInput {
	return CreatePostInput{
		UserID:  strings.TrimSpace(input.UserID),
		Title:   strings.TrimSpace(input.Title),
		Content: strings.TrimSpace(input.Content),
	}
}

func ValidateInput(input CreatePostInput) error {
	validate := validator.New()
	return validate.Struct(input)
}

func UsecaseValidate(input CreatePostInput) error {
	validate := validator.New()
	return validate.Struct(input)
}

func ValidationToMap(err error) fiber.Map {
	mapped := fiber.Map{}
	if err == nil {
		return mapped
	}
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, verr := range verrs {
			switch strings.ToLower(verr.Field()) {
			case "userid":
				mapped["ErrorUserID"] = postFriendlyMessage(verr)
			case "title":
				mapped["ErrorTitle"] = postFriendlyMessage(verr)
			case "content":
				mapped["ErrorContent"] = postFriendlyMessage(verr)
			}
		}
		return mapped
	}
	mapped["ErrorGeneral"] = err.Error()
	return mapped
}

func postFriendlyMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required."
	case "min":
		return "Must be at least 2 characters."
	case "max":
		return "Too long."
	default:
		return "Invalid value."
	}
}

func baseData(c fiber.Ctx, data fiber.Map) fiber.Map {
	var store *session.Store
	flashPop := middleware.PopFlash(store, c)
	flash := middleware.ConsumeFlash(store, c)
	csrfToken := csrf.TokenFromContext(c)
	base := fiber.Map{
		"FlashSuccess": flash.Success,
		"FlashError":   flash.Error,
		"CsrfToken":    csrfToken,
		"FlashType":    flashPop.Type,
		"FlashMessage": flashPop.Message,
	}
	for k, v := range data {
		base[k] = v
	}
	return base
}

func BindInput(c fiber.Ctx) (CreatePostInput, fiber.Map, error) {
	input := NormalizeInput(CreatePostInput{

		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
	})
	if err := UsecaseValidate(input); err != nil {
		return input, ValidationToMap(err), err
	}
	return input, fiber.Map{}, nil
}

func RenderCreateWithErrors(c fiber.Ctx, input CreatePostInput, data fiber.Map, err error) error {
	base := baseData(c, fiber.Map{
		"PageTitle":   "Create post",
		"FormAction":  "/web/post/store",
		"SubmitLabel": "Create post",
		"FormMode":    "create",
		"Post":        input,
	})
	for k, v := range data {
		base[k] = v
	}
	if len(data) == 0 {
		base["ErrorGeneral"] = err.Error()
	}
	return c.Status(http.StatusUnprocessableEntity).Render("posts/create", base)
}

func RenderEditWithErrors(c fiber.Ctx, id string, input CreatePostInput, data fiber.Map, err error) error {
	base := baseData(c, fiber.Map{
		"PageTitle":   "Edit post",
		"FormAction":  fmt.Sprintf("/web/post/%s/update", id),
		"SubmitLabel": "Update post",
		"FormMode":    "edit",
		"PostID":      id,
		"Post":        input,
	})
	for k, v := range data {
		base[k] = v
	}
	if len(data) == 0 {
		base["ErrorGeneral"] = err.Error()
	}
	return c.Status(http.StatusUnprocessableEntity).Render("posts/edit", base)
}
