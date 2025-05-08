package editing_events

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func SaveEvent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат данных"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Валидация данных
	if event.Description == "" || event.SportCategory == "" || event.SportType == "" || event.PriceStatus == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Все поля обязательны для заполнения"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Duration <= 0 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Продолжительность должна быть больше 0"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Если event_id указан, это обновление существующего ивента
	if event.EventID == -1 {
		if err := database.DB.
			Table(`public."Event"`).
			Where("event_id = ?", event.EventID).
			Updates(&event).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обновления ивента"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// Иначе создаём новый ивент
		if err := database.DB.
			Table(`public."Event"`).
			Create(&event).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка создания ивента"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Формируем ответ
	response := map[string]any{
		"success":  true,
		"event_id": event.EventID,
		"message":  "Ивент успешно сохранён",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
