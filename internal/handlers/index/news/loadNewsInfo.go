package news

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

func LoadNewsInfo(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	newsIdStr := r.URL.Query().Get("news_id")
	newsId, err := strconv.Atoi(newsIdStr)
	if err != nil || newsId < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный ID новости"})
		return
	}

	_, userId, userRights, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	isModerator := userRights%16 >= 8 || userRights/2147483648 == 1
	if !isModerator {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Доступ запрещён"})
		return
	}

	var news models.News
	query := database.DB.Where("news_id = ?", newsId)

	if userRights/2147483648 != 1 {
		query = query.Where("creator_id = ?", userId)
	}
	if err := query.First(&news).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Новость не найдена"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Success bool        `json:"success"`
		News    models.News `json:"news"`
	}{
		Success: true,
		News:    news,
	})
}
