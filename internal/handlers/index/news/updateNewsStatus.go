package news

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func UpdateNewsStatus(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	var requestData struct {
		NewsID int    `json:"news_id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат данных"})
		return
	}

	if requestData.NewsID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный news_id"})
		return
	}

	if requestData.Status != "ACTIVE" && requestData.Status != "ARCHIVE" && requestData.Status != "DRAFT" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный статус новости"})
		return
	}

	var news models.News
	if err := database.DB.Where("news_id = ?", requestData.NewsID).First(&news).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Новость не найдена"})
		return
	}

	news.Status = requestData.Status

	if err := database.DB.Save(&news).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при обновлении статуса новости"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Статус новости успешно обновлён"})
}
