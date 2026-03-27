package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/jiehui555/meaw-oa/internal/model"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Create(&model.User{
		Name:     "admin",
		Phone:    "00000000000",
		Email:    "admin@meaw.com",
		Password: string(hashed),
	})

	return db
}

func setupApp(t *testing.T, db *gorm.DB) *fiber.App {
	t.Helper()

	app := fiber.New()
	api := app.Group("/api")
	userHandler := NewUserHandler(db)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.Refresh)

	return app
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func doRequest(t *testing.T, app *fiber.App, method, path string, body any) (*http.Response, apiResponse) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var res apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return resp, res
}
