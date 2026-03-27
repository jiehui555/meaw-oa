package handler

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	app := setupApp(t, db)

	t.Run("success", func(t *testing.T) {
		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":     "admin",
			"password": "password",
		})

		if res.Code != 0 {
			t.Errorf("expected code 0, got %d: %s", res.Code, res.Message)
		}

		var data map[string]string
		if err := json.Unmarshal(res.Data, &data); err != nil {
			t.Fatalf("failed to unmarshal data: %v", err)
		}
		if data["token"] == "" {
			t.Error("expected token to be non-empty")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		resp, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":     "admin",
			"password": "wrong",
		})

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if res.Code != 401 {
			t.Errorf("expected code 401, got %d", res.Code)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		_, res := doRequest(t, app, "POST", "/api/login", map[string]string{
			"name":     "nonexistent",
			"password": "password",
		})

		if res.Code != 401 {
			t.Errorf("expected code 401, got %d", res.Code)
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
