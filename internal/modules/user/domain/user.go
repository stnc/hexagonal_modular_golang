package domain

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gopkg.in/go-playground/validator.v9"
	"hexagonalapp/internal/modules/user/adapters/inbound/http/middleware"
	"net/http"
	"strings"
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CreateUserInput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"  validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email,max=150"`
}

type UpdateUserInput struct {
	ID    uint   `validate:"required,gt=0"`
	Name  string `validate:"required,min=2,max=100"`
	Email string `validate:"required,email,max=150"`
}

var ErrEmailAlreadyUsed = errors.New("email already used")

func ValidateInput(input CreateUserInput) error {
	validate := validator.New()

	if err := validate.Struct(input); err != nil {
		return err
	}
	return nil
}

func UsecaseValidate(input CreateUserInput) error {
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
			case "name":
				mapped["ErrorName"] = userFriendlyMessage(verr)
			case "email":
				mapped["ErrorEmail"] = userFriendlyMessage(verr)
			case "age":
				mapped["ErrorAge"] = userFriendlyMessage(verr)
			}
		}
		return mapped
	}
	mapped["ErrorGeneral"] = err.Error()
	return mapped
}

func userFriendlyMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required."
	case "email":
		return "Please enter a valid email address."
	case "min":
		return "Must be at least 3 characters."
	case "max":
		return "Too long."
	case "gte":
		return "Must be 13 or older."
	case "lte":
		return "Must be 120 or younger."
	default:
		return "Invalid value."
	}
}

func baseData(c fiber.Ctx, data fiber.Map) fiber.Map {
	var store *session.Store
	flash := middleware.PopFlash(store, c)
	csrfToken := csrf.TokenFromContext(c)
	base := fiber.Map{
		"FlashSuccess": flash.Success,
		// "FlashError":   message,// TODO: bunu bak
		"CsrfToken":    csrfToken,
		"FlashType":    flash.Type,
		"FlashMessage": flash.Message,
	}
	for k, v := range data {
		base[k] = v
	}
	return base
}

func BindInput(c fiber.Ctx) (CreateUserInput, fiber.Map, error) {
	input := CreateUserInput{
		Name:  strings.TrimSpace(c.FormValue("name")),
		Email: strings.TrimSpace(c.FormValue("email")),
	}

	if err := UsecaseValidate(input); err != nil {
		return input, ValidationToMap(err), err
	}
	return input, fiber.Map{}, nil
}





func RenderCreateWithErrors(c fiber.Ctx, input CreateUserInput, data fiber.Map, err error) error {
	base := baseData(c, fiber.Map{
		"PageTitle":   "Create user",
		"FormAction":  "/users",
		"SubmitLabel": "Create user",
		"FormMode":    "create",
		"User":        input,
		// "Age":         strconv.Itoa(input.Age),
	})
	for k, v := range data {
		base[k] = v
	}
	if len(data) == 0 {
		base["ErrorGeneral"] = err.Error()
	}
	return c.Status(http.StatusUnprocessableEntity).Render("users/create", base)
}

