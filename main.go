package main

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"github.com/jiehui555/meaw-oa/internal/common"
	"github.com/jiehui555/meaw-oa/internal/config"
	"github.com/jiehui555/meaw-oa/internal/database"
	"github.com/jiehui555/meaw-oa/internal/handler"
	"github.com/jiehui555/meaw-oa/internal/logger"
	"github.com/jiehui555/meaw-oa/internal/middleware"
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

	admin := api.Group("/admin", middleware.Auth(db), middleware.Admin())
	admin.Get("/dashboard", func(c fiber.Ctx) error {
		return common.Success(c, "管理员仪表板")
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("服务器启动中", "addr", addr)
	if err := app.Listen(addr); err != nil {
		slog.Error("服务器启动失败", "error", err)
		panic(err)
	}
}
