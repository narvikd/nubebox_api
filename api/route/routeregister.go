package route

import (
	"api/api/jsonresponse"
	"api/db/dbengine"
	"api/internal/app"
	"api/internal/cfg"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

type ApiCtx struct {
	HttpServer *fiber.App
	Config     *cfg.Config
	DB         *sql.DB
	Query      *dbengine.Queries
}

func newRouteCtx(app *app.App) *ApiCtx {
	c := &ApiCtx{
		HttpServer: app.HttpServer,
		Config:     app.Config,
		DB:         app.DB,
		Query:      app.Query,
	}
	return c
}

func Register(app *app.App) {
	routes(app.HttpServer, newRouteCtx(app))
}

func routes(app *fiber.App, route *ApiCtx) {
	api := app.Group("/api/v1")
	api.Get("/", route.getRoot)
}

func (c *ApiCtx) getRoot(fiberCtx *fiber.Ctx) error {
	return jsonresponse.OK(fiberCtx, "Hello")
}
