package routers

import (
	"api_v2/database"
	"api_v2/internal/middlewares"
	"api_v2/internal/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetUpRouter() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Get("/insoske", monitor.New()) //For monitoring

	const ROUTE_PREFIX = "api/v1/web"

	var (
		userService    = users.NewService(database.DB)
		userController = users.NewController(userService)
	)

	public := app.Group(ROUTE_PREFIX + "/users")
	{
		public.Post("/register", userController.Register)
		public.Post("/login", userController.Login)
	}

	protected := app.Group(ROUTE_PREFIX + "/users")
	protected.Use(middlewares.Middleware())
	{
		protected.Get("/", userController.GetUsers)
	}

	return app

}
