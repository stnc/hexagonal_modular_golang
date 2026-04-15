package app

import (
		"time"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
		"github.com/gofiber/fiber/v3/extractors"
			"github.com/gofiber/fiber/v3/middleware/session"

	"github.com/gofiber/fiber/v3/middleware/csrf"
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

func (h *Handlers) Run(app *fiber.App, store *session.Store, cfgEnv string) {


	app.Use(csrf.New(csrf.Config{
		Session:        store,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		CookieSecure:   cfgEnv == "PRODUCTION",
		Extractor:      extractors.FromForm("_csrf"),
		IdleTimeout:    30 * time.Minute,
		CookieName:     "csrf_",
	}))

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
