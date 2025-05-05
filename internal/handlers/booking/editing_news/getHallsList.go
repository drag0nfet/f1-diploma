package editing_news

import (
	"diploma/internal/database"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func GetHallsList(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var halls []struct {
		HallID int    `json:"hall_id"`
		Name   string `json:"name"`
	}

	if err := database.DB.
		Table("\"Hall\"").
		Select("hall_id, name").
		Order("hall_id ASC").
		Find(&halls).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки залов"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"halls":   halls,
	})
}
