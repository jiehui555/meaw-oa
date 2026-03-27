package database

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
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
		log.Fatal(fmt.Errorf("failed to hash password: %w", err))
	}

	admin := model.User{
		Name:     "admin",
		Phone:    "00000000000",
		Email:    "admin@meaw.com",
		Password: string(hashed),
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatal(fmt.Errorf("failed to seed admin user: %w", err))
	}

	log.Println("Super admin user created (admin / password)")
}
