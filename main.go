package main

import (
	"api/api/middleware"
	"api/api/route"
	"api/internal/app"
	"github.com/gofiber/fiber/v2"
	"github.com/narvikd/errorskit"
)

func main() {
	const apiAddr = "0.0.0.0:3001"
	if errListen := newApi().Listen(apiAddr); errListen != nil {
		errorskit.FatalWrap(errListen, "api can't be started")
	}
}

func newApi() *fiber.App {
	a := app.NewApp()
	middleware.InitMiddlewares(a.HttpServer)
	route.Register(a)
	return a.HttpServer
}
