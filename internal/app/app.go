package app

import (
	"api/db/dbconfig"
	"api/db/dbengine"
	"api/internal/cfg"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/narvikd/fiberparser"
	"log"
	"time"
)

// App is a simple struct to include a collection of tools that the application could need to operate.
//
// For example the pointer to DB that is used all over the application.
//
// This way the application can avoid the use of global variables.
type App struct {
	HttpServer *fiber.App
	Config     *cfg.Config
	DB         *sql.DB
	Query      *dbengine.Queries
}

func NewApp() *App {
	config := cfg.InitCfg()

	serv := fiber.New(fiber.Config{
		AppName:      "NubeBox",
		IdleTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return fiberparser.RegisterErrorHandler(ctx, err)
		},
		BodyLimit: 10 * 1024 * 1024 * 1024, // In GB
	})

	db, errDBInit := dbconfig.InitDB(config)
	if errDBInit != nil {
		log.Fatalln(errDBInit)
	}

	a := &App{
		HttpServer: serv,
		Config:     config,
		DB:         db,
		Query:      dbengine.New(db),
	}
	return a
}
