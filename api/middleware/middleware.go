package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func InitMiddlewares(app *fiber.App) {
	initLoggerMW(app)
	initCorsMW(app)
	initRecoverMw(app)
}

func initCorsMW(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowCredentials: true,
		}),
	)
}

func initRecoverMw(app *fiber.App) {
	app.Use(recover.New())
}

func initLoggerMW(app *fiber.App) {
	app.Use(
		logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		}),
	)
}
