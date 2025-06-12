package news

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DeleteNews(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	_, _, userRights, response := services.CheckAuthCookie(r)
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

	vars := mux.Vars(r)
	newsIdStr, ok := vars["news_id"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "ID новости не указан"})
		return
	}
	newsId, err := strconv.ParseInt(newsIdStr, 10, 64)
	if err != nil || newsId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный ID новости"})
		return
	}

	var news models.News
	if err := database.DB.Where("news_id = ?", newsId).First(&news).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Новость не найдена"})
		return
	}

	if err := database.DB.Delete(&news).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Ошибка при удалении новости",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Новость успешно удалена"})
}
