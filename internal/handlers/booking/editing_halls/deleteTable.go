package editing_halls

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tableIDParam := r.URL.Query().Get("table_id")
	if tableIDParam == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указан table_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tableID, err := strconv.Atoi(tableIDParam)
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный table_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = database.DB.
		Table(`public."Table"`).
		Where("table_id = ?", tableID).
		Delete(&models.Table{}).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка удаления стола"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Стол удалён"})
}
