package database

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/jiehui555/meaw-oa/internal/model"
)

func Init(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect database: %w", err))
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(fmt.Errorf("failed to migrate database: %w", err))
	}

	return db
}
