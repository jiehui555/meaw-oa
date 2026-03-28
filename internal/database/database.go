package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/jiehui555/meaw-oa/internal/model"
)

// slogLogger GORM 日志适配器，使用 slog 输出日志
type slogLogger struct {
	LogLevel logger.LogLevel
}

func newSlogLogger() logger.Interface {
	return &slogLogger{
		LogLevel: logger.Info,
	}
}

func (l *slogLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *slogLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		slog.Info(fmt.Sprintf(msg, data...))
	}
}

func (l *slogLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		slog.Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *slogLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		slog.Error(fmt.Sprintf(msg, data...))
	}
}

func (l *slogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	attrs := []slog.Attr{
		slog.String("file", utils.FileWithLineNum()),
		slog.String("elapsed", elapsed.String()),
		slog.Int64("rows", rows),
		slog.String("sql", sql),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		attrs = append(attrs, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "数据库", attrs...)
	case elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn:
		slog.LogAttrs(ctx, slog.LevelWarn, "慢查询", attrs...)
	case l.LogLevel >= logger.Info:
		slog.LogAttrs(ctx, slog.LevelInfo, "数据库", attrs...)
	}
}

// Init 初始化数据库连接并执行迁移
// dbPath: SQLite 数据库文件路径
func Init(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newSlogLogger(),
	})
	if err != nil {
		slog.Error("连接数据库失败", "error", err)
		panic(err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Captcha{}); err != nil {
		slog.Error("数据库迁移失败", "error", err)
		panic(err)
	}

	seedAdmin(db)

	return db
}

// seedAdmin 初始化默认管理员账户
func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Where("name = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("密码加密失败", "error", err)
		panic(err)
	}

	admin := model.User{
		Name:     "admin",
		Phone:    "00000000000",
		Email:    "admin@meaw.com",
		Password: string(hashed),
	}

	if err := db.Create(&admin).Error; err != nil {
		slog.Error("初始化管理员用户失败", "error", err)
		panic(err)
	}

	slog.Info("超级管理员用户已创建", "name", "admin", "password", "password")
}
