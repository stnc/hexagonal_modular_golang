package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	postsmongo "hexagonalapp/internal/modules/posts/adapters/outbound/mongodb"
	postspg "hexagonalapp/internal/modules/posts/adapters/outbound/postgres"
	postsredis "hexagonalapp/internal/modules/posts/adapters/outbound/redis"
	postsapp "hexagonalapp/internal/modules/posts/app"
	usermongo "hexagonalapp/internal/modules/user/adapters/outbound/mongodb"
	userpg "hexagonalapp/internal/modules/user/adapters/outbound/postgres"
	userredis "hexagonalapp/internal/modules/user/adapters/outbound/redis"
	userapp "hexagonalapp/internal/modules/user/app"
	redisplatform "hexagonalapp/internal/platform/cache/redis"
	"hexagonalapp/internal/platform/config"
	mongoplatform "hexagonalapp/internal/platform/database/mongodb"
	postgresplatform "hexagonalapp/internal/platform/database/postgres"
	api "hexagonalapp/internal/transport/api"
	web "hexagonalapp/internal/transport/web"
	"time"
	// postshttp "hexagonalapp/internal/modules/posts/adapters/inbound/http"
	// userhttp "hexagonalapp/internal/modules/user/adapters/inbound/http"
)

func Runner(app *fiber.App) error {
	var store *session.Store
	cfg := config.Load()

	pgDB := postgresplatform.DbConnect(cfg)

	MongoDBURI := "mongodb://" + cfg.MongoDBUser + ":" + cfg.MongoDBPassword + "@localhost:27017/?authSource=admin"

	mongoClient, err := mongoplatform.Open(context.Background(), MongoDBURI)
	if err != nil {
		return err
	}

	rdb := redisplatform.New(cfg.RedisHost+":"+cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)

	userRepo := userpg.New(pgDB)
	if err := userRepo.AutoMigrate(); err != nil {
		return err
	}
	userCache := userredis.New(rdb)
	userAudit := usermongo.New(mongoClient.Database(cfg.MongoDBName).Collection("user_events"))
	userService := userapp.New(userRepo, userCache, userAudit)
	// userHandler := userhttp.NewAPI(userService)

	postsRepo := postsmongo.New(mongoClient.Database(cfg.MongoDBName).Collection("posts"))
	postsCache := postsredis.New(rdb)
	postsMirror := postspg.New(pgDB)
	if err := postsMirror.AutoMigrate(); err != nil {
		return err
	}
	postsService := postsapp.New(postsRepo, postsCache, postsMirror)
	//postsHandler := postshttp.NewAPI(postsService)

	store = session.NewStore(session.Config{
		IdleTimeout:       30 * time.Minute,
		AbsoluteTimeout:   24 * time.Hour, // Force expire after 24 hours regardless of activity
		CookieHTTPOnly:    true,
		CookieSameSite:    "Lax",
		CookieSecure:      cfg.EnvName == "PRODUCTION",
		CookieSessionOnly: false,
	
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(session.New(session.Config{Store: store}))

	/*
		apiRunning := api.Handlers{
				UserHandler:  userhttp.NewAPI(userService),
				PostsHandler: postshttp.NewAPI(postsService),
			}
					apiRunning.Run(app)
	*/
	// app.Use(NotFound) // 404 // TODO : how its work ?

	switch cfg.App {
	case "api":
		apiWORker := api.NewHandlers(userService, postsService)
		apiWORker.Run(app)
	case "web":
		webWORker := web.NewHandlers(userService, postsService)
		webWORker.Run(app, store, cfg.EnvName)

	default:
		apiWORker := api.NewHandlers(userService, postsService)
		apiWORker.Run(app)
		webWORker := web.NewHandlers(userService, postsService)
		webWORker.Run(app, store, cfg.EnvName)
	}

	return app.Listen(fmt.Sprintf(":%s", "9999"))

}
func NotFound(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render("404", nil)
}
