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
		slog.LogAttrs(ctx, slog.LevelError, "gorm", attrs...)
	case elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn:
		slog.LogAttrs(ctx, slog.LevelWarn, "slow query", attrs...)
	case l.LogLevel >= logger.Info:
		slog.LogAttrs(ctx, slog.LevelInfo, "gorm", attrs...)
	}
}

func Init(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newSlogLogger(),
	})
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		panic(err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		slog.Error("failed to migrate database", "error", err)
		panic(err)
	}

	seedAdmin(db)

	return db
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Where("name = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		panic(err)
	}

	admin := model.User{
		Name:     "admin",
		Phone:    "00000000000",
		Email:    "admin@meaw.com",
		Password: string(hashed),
	}

	if err := db.Create(&admin).Error; err != nil {
		slog.Error("failed to seed admin user", "error", err)
		panic(err)
	}

	slog.Info("Super admin user created", "name", "admin", "password", "password")
}
