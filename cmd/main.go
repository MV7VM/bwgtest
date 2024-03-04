package main

import (
	"bwg/config"
	"bwg/internal/handlers"
	"bwg/internal/repository.go"
	"bwg/internal/service"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
	"time"
)

const port = ":3000"

type App struct {
	routers    *fiber.App
	repository *repository_go.Repository
	service    *service.Service
	handlers   *handlers.Handlers
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
	app := &App{}
	app.routers = fiber.New()
	cfg := config.Config_load()
	app.repository = repository_go.New(cfg)
	go func() {
		for {
			err := app.repository.CurrentPrice()
			if err != nil {
				slog.Error("Error in Updating price:", err)
			}
			time.Sleep(time.Minute)
		}
	}()
	app.service = service.New(app.repository)
	app.handlers = handlers.New(app.service)
	app.routers.Post("/add_ticker", app.handlers.Post)
	app.routers.Post("/fetch", app.handlers.Get)
	err := app.routers.Listen(":3000")
	if err != nil {
		logger.Error("Cann't listen ", port, "\n", err)
	}
}
