package main

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"github.com/jiehui555/meaw-oa/internal/config"
	"github.com/jiehui555/meaw-oa/internal/database"
	"github.com/jiehui555/meaw-oa/internal/logger"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogPath)

	db := database.Init(cfg.DBPath)
	_ = db

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("server starting", "addr", addr)
	if err := app.Listen(addr); err != nil {
		slog.Error("server failed to start", "error", err)
		panic(err)
	}
}
