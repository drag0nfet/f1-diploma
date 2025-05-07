package editing_halls

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetHall(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
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

	var hall models.Hall
	if err = database.DB.
		Table(`public."Hall"`).
		Where("hall_id = ?", hallID).
		First(&hall).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Зал не найден"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Загрузка фотографий
	var photos []models.HallPhoto
	if err = database.DB.
		Table(`public."HallPhotos"`).
		Where("hall_id = ?", hallID).
		Order("created_at ASC").
		Find(&photos).Error; err != nil {
		log.Printf("Ошибка загрузки фотографий для зала %d: %v", hallID, err)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки фотографий"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tables []models.Table
	if err = database.DB.
		Table(`public."Table"`).
		Where("hall_id = ?", hallID).
		Order("seats DESC, table_name ASC").
		Find(&tables).Error; err != nil {
		log.Printf("Ошибка загрузки столов для зала %d: %v", hallID, err)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки столов"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := map[string]any{
		"success": true,
		"hall": map[string]any{
			"hall_id":     hall.HallID,
			"name":        hall.Name,
			"description": hall.Description,
			"album":       photos,
			"tables":      tables,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования ответа для зала %d: %v", hallID, err)
	}
}
