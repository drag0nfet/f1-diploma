package editing_halls

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

func DeleteHall(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	hallIDParam := r.URL.Query().Get("hall_id")
	if hallIDParam == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указан hall_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hallID, err := strconv.Atoi(hallIDParam)
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный hall_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := database.DB.
		Table(`public."Hall"`).
		Where("hall_id = ?", hallID).
		Delete(&models.Hall{})
	if result.Error != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка удаления зала"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Зал не найден"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Зал удалён"})
}
