package logger

import (
	"log/slog"
	"os"
)

// Init 初始化日志系统
// logPath: 日志文件路径
func Init(logPath string) {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("打开日志文件失败：" + err.Error())
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}
