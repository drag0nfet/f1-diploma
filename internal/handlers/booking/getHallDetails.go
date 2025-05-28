package booking

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func GetHallDetails(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем hallId из пути
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат URL"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hallIdStr := pathParts[3]
	hallId, err := strconv.Atoi(hallIdStr)
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный hall_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получаем информацию о зале
	var hall models.Hall
	if err := database.DB.Table(`public."Hall"`).Where("hall_id = ?", hallId).First(&hall).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки зала"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получаем фотографии зала
	var photos []models.HallPhoto
	if err := database.DB.Table(`public."HallPhotos"`).Where("hall_id = ?", hallId).Find(&photos).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки фотографий"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"success": true,
		"hall":    hall,
		"photos":  photos,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
