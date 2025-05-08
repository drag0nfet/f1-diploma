package editing_events

import (
	"diploma/internal/database"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func GetEventsList(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var events []struct {
		EventID     int    `json:"event_id"`
		Description string `json:"description"`
	}

	if err := database.DB.
		Table("\"Event\"").
		Select("event_id, description").
		Order("event_id ASC").
		Find(&events).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки ивентов"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"events":  events,
	})
}
