package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	postshttp "hexagonalapp/internal/modules/posts/adapters/inbound/http"
	postsapp "hexagonalapp/internal/modules/posts/app"
	userhttp "hexagonalapp/internal/modules/user/adapters/inbound/http"
	userapp "hexagonalapp/internal/modules/user/app"
)

type Handlers struct {
	UserHandler  *userhttp.WebHandler
	PostsHandler *postshttp.WEBHandler
}

func NewHandlers(
	userService *userapp.Service,
	postsService *postsapp.Service,
) *Handlers {
	return &Handlers{
		UserHandler:  userhttp.NewWEB(userService),
		PostsHandler: postshttp.NewWEB(postsService),
	}
}

func (h *Handlers) Run(app *fiber.App) {


	app.Get("/", func(c fiber.Ctx) error {
		// Render with and extends
		return c.Render("homepage", fiber.Map{
			"Title": "this homepage",
		})
	})

	app.Get("/dashboard", func(c fiber.Ctx) error {
		// Render with and extend
		return c.Render("homepage", fiber.Map{
			"Title": "this homepage",
		})
	})



	app.Get("/public*", static.New("./public"))

	web := app.Group("/web")
	h.UserHandler.Register(web)
	h.PostsHandler.Register(web)

}
