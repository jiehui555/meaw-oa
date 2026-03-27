package main

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"github.com/jiehui555/meaw-oa/internal/config"
	"github.com/jiehui555/meaw-oa/internal/database"
	"github.com/jiehui555/meaw-oa/internal/handler"
	"github.com/jiehui555/meaw-oa/internal/logger"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogPath)

	db := database.Init(cfg.DBPath)

	app := fiber.New()

	api := app.Group("/api")
	userHandler := handler.NewUserHandler(db)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.Refresh)

	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("server starting", "addr", addr)
	if err := app.Listen(addr); err != nil {
		slog.Error("server failed to start", "error", err)
		panic(err)
	}
}
