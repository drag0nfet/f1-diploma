package editing_events

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetEvent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	eventIDParam := r.URL.Query().Get("event_id")
	if eventIDParam == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указан event_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный event_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var event models.Event
	if err = database.DB.
		Table(`public."Event"`).
		Where("event_id = ?", eventID).
		First(&event).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ивент не найден"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Формируем ответ
	response := map[string]any{
		"success": true,
		"event":   event,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования ответа для ивента %d: %v", eventID, err)
	}
}
