package index

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func LoadNew(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	newsIdStr, ok := vars["newsId"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "ID новости не указан"})
		return
	}
	newsID, err := strconv.ParseInt(newsIdStr, 10, 64)
	if err != nil || newsID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный ID новости"})
		return
	}

	var news models.News
	if err = database.DB.Where("news_id = ?", newsID).First(&news).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Новость не найдено"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		News    models.News
		Success bool
	}{
		News:    news,
		Success: true,
	})
}
