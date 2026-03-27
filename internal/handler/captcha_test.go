package handler

import (
	"encoding/json"
	"testing"

	"github.com/jiehui555/meaw-oa/internal/model"
)

func TestGetCaptcha(t *testing.T) {
	db := setupTestDB(t)
	app := setupApp(t, db)

	t.Run("success", func(t *testing.T) {
		_, res := doRequest(t, app, "GET", "/api/captcha", nil)

		if res.Code != 0 {
			t.Errorf("expected code 0, got %d: %s", res.Code, res.Message)
		}

		var data map[string]string
		if err := json.Unmarshal(res.Data, &data); err != nil {
			t.Fatalf("failed to unmarshal data: %v", err)
		}

		if data["captcha_id"] == "" {
			t.Error("expected captcha_id to be non-empty")
		}
		if data["captcha_img"] == "" {
			t.Error("expected captcha_img to be non-empty")
		}

		var count int64
		db.Model(&model.Captcha{}).Where("captcha_id = ?", data["captcha_id"]).Count(&count)
		if count != 1 {
			t.Errorf("expected captcha record in db, got %d", count)
		}
	})
}
