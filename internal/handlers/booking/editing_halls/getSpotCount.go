package editing_halls

import (
	"diploma/internal/database"
	"diploma/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetSpotCount(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tableIDParam := r.URL.Query().Get("table_id")
	if tableIDParam == "" {
		log.Printf("Отсутствует параметр table_id")
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указан table_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tableID, err := strconv.Atoi(tableIDParam)
	if err != nil {
		log.Printf("Некорректный table_id: %v", tableIDParam)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный table_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int64
	if err = database.DB.
		Table(`public."Spot"`).
		Where("table_id = ?", tableID).
		Count(&count).Error; err != nil {
		log.Printf("Ошибка подсчёта спотов для стола %d: %v", tableID, err)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка подсчёта спотов"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"count":   count,
	})
}
