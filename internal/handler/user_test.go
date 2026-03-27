package handler

import (
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/jiehui555/meaw-oa/internal/model"
	"gorm.io/gorm"
)

func getCaptcha(t *testing.T, app *fiber.App, db *gorm.DB) (captchaID, answer string) {
	t.Helper()

	_, res := doRequest(t, app, "GET", "/api/captcha", nil)
	var data map[string]string
	json.Unmarshal(res.Data, &data)
	captchaID = data["captcha_id"]

	var captcha model.Captcha
	db.Where("captcha_id = ?", captchaID).First(&captcha)
	answer = captcha.Answer
	return
}

func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	app := setupApp(t, db)

	t.Run("success", func(t *testing.T) {
		captchaID, answer := getCaptcha(t, app, db)

		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":           "admin",
			"password":       "password",
			"captcha_id":     captchaID,
			"captcha_answer": answer,
		})

		if res.Code != 0 {
			t.Errorf("expected code 0, got %d: %s", res.Code, res.Message)
		}

		var data map[string]string
		if err := json.Unmarshal(res.Data, &data); err != nil {
			t.Fatalf("failed to unmarshal data: %v", err)
		}
		if data["access_token"] == "" {
			t.Error("expected access_token to be non-empty")
		}
		if data["refresh_token"] == "" {
			t.Error("expected refresh_token to be non-empty")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		captchaID, answer := getCaptcha(t, app, db)

		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":           "admin",
			"password":       "wrong",
			"captcha_id":     captchaID,
			"captcha_answer": answer,
		})

		if res.Code != 401 {
			t.Errorf("expected code 401, got %d", res.Code)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		captchaID, answer := getCaptcha(t, app, db)

		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":           "nonexistent",
			"password":       "password",
			"captcha_id":     captchaID,
			"captcha_answer": answer,
		})

		if res.Code != 401 {
			t.Errorf("expected code 401, got %d", res.Code)
		}
	})

	t.Run("wrong captcha", func(t *testing.T) {
		captchaID, _ := getCaptcha(t, app, db)

		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":           "admin",
			"password":       "password",
			"captcha_id":     captchaID,
			"captcha_answer": "0000",
		})

		if res.Code != 400 {
			t.Errorf("expected code 400, got %d", res.Code)
		}
	})

	t.Run("empty fields", func(t *testing.T) {
		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name": "",
		})

		if res.Code != 400 {
			t.Errorf("expected code 400, got %d", res.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		_, res := doRequest(t, app, "POST", "/api/login", "not json")

		if res.Code != 400 {
			t.Errorf("expected code 400, got %d", res.Code)
		}
	})
}

func TestRefresh(t *testing.T) {
	db := setupTestDB(t)
	app := setupApp(t, db)

	getRefreshToken := func(t *testing.T) string {
		t.Helper()
		captchaID, answer := getCaptcha(t, app, db)
		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":           "admin",
			"password":       "password",
			"captcha_id":     captchaID,
			"captcha_answer": answer,
		})
		var data map[string]string
		json.Unmarshal(res.Data, &data)
		return data["refresh_token"]
	}

	t.Run("success", func(t *testing.T) {
		refreshToken := getRefreshToken(t)

		_, res := doRequest(t, app, "POST", "/api/refresh", map[string]string{
			"refresh_token": refreshToken,
		})

		if res.Code != 0 {
			t.Errorf("expected code 0, got %d: %s", res.Code, res.Message)
		}

		var data map[string]string
		if err := json.Unmarshal(res.Data, &data); err != nil {
			t.Fatalf("failed to unmarshal data: %v", err)
		}
		if data["access_token"] == "" {
			t.Error("expected access_token to be non-empty")
		}
		if data["refresh_token"] == "" {
			t.Error("expected refresh_token to be non-empty")
		}
	})

	t.Run("use access token to refresh should fail", func(t *testing.T) {
		captchaID, answer := getCaptcha(t, app, db)

		_, loginRes := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":           "admin",
			"password":       "password",
			"captcha_id":     captchaID,
			"captcha_answer": answer,
		})
		var loginData map[string]string
		json.Unmarshal(loginRes.Data, &loginData)

		_, res := doRequest(t, app, "POST", "/api/refresh", map[string]string{
			"refresh_token": loginData["access_token"],
		})

		if res.Code != 401 {
			t.Errorf("expected code 401, got %d", res.Code)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		_, res := doRequest(t, app, "POST", "/api/refresh", map[string]string{})

		if res.Code != 400 {
			t.Errorf("expected code 400, got %d", res.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, res := doRequest(t, app, "POST", "/api/refresh", map[string]string{
			"refresh_token": "invalid.token.here",
		})

		if res.Code != 401 {
			t.Errorf("expected code 401, got %d", res.Code)
		}
	})
}
