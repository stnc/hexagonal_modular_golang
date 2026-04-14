package app

import (
	"github.com/gofiber/fiber/v3"

	postshttp "hexagonalapp/internal/modules/posts/adapters/inbound/http"
	postsapp "hexagonalapp/internal/modules/posts/app"
	userhttp "hexagonalapp/internal/modules/user/adapters/inbound/http"
	userapp "hexagonalapp/internal/modules/user/app"
)

type Handlers struct {
	UserHandler  *userhttp.APIHandler
	PostsHandler *postshttp.APIHandler
}

func NewHandlers(
	userService *userapp.Service,
	postsService *postsapp.Service,
) *Handlers {
	return &Handlers{
		UserHandler:  userhttp.NewAPI(userService),
		PostsHandler: postshttp.NewAPI(postsService),
	}
}

func (h *Handlers) Run(app *fiber.App) {

	api := app.Group("/api")
	h.UserHandler.Register(api)
	h.PostsHandler.Register(api)

}
