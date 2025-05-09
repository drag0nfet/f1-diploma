package editing_events

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodDelete {
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

	result := database.DB.
		Table(`public."Event"`).
		Where("event_id = ?", eventID).
		Delete(&models.Event{})
	if result.Error != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка удаления ивента"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ивент не найден"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Ивент удалён"})
}
