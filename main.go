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
	// 加载配置文件
	cfg := config.Load()

	// 初始化日志系统
	logger.Init(cfg.LogPath)

	// 初始化数据库连接
	db := database.Init(cfg.DBPath)

	// 创建 Fiber 应用实例
	app := fiber.New()

	// 注册 API 路由组
	api := app.Group("/api")
	userHandler := handler.NewUserHandler(db)
	captchaHandler := handler.NewCaptchaHandler(db)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.Refresh)
	api.Get("/captcha", captchaHandler.GetC)

	// 注册管理员路由组，需要认证和管理员权限
	admin := api.Group("/admin", middleware.Auth(db), middleware.Admin())
	admin.Get("/dashboard", func(c fiber.Ctx) error {
		return common.Success(c, "管理员仪表板")
	})

	// 启动 HTTP 服务器
	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("服务器启动中", "addr", addr)
	if err := app.Listen(addr); err != nil {
		slog.Error("服务器启动失败", "error", err)
		panic(err)
	}
}
