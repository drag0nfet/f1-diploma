package news

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

func LoadNews(w http.ResponseWriter, r *http.Request) {
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

	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if status != "ACTIVE" && status != "ARCHIVE" && status != "DRAFT" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный статус"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный номер страницы"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректное количество записей"})
		return
	}

	_, userId, userRights, _ := services.CheckAuthCookie(r)
	if (status == "DRAFT" || status == "ARCHIVE") && userRights%16 < 8 && userRights/2147483648 != 1 {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Доступ запрещён"})
		return
	}

	var total int64
	query := database.DB.Model(&models.News{}).Where("status = ?", status)
	if status == "DRAFT" || status == "ARCHIVE" {
		query = query.Where("creator_id = ?", userId)
	}
	if err := query.Count(&total).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка подсчёта новостей"})
		return
	}

	var news []models.News
	offset := (page - 1) * limit
	query = database.DB.Where("status = ?", status)
	if status == "DRAFT" || status == "ARCHIVE" {
		query = query.Where("creator_id = ?", userId)
	}
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&news).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки новостей"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Success bool          `json:"success"`
		AllNews []models.News `json:"all_news"`
		Total   int64         `json:"total"`
	}{
		Success: true,
		AllNews: news,
		Total:   total,
	})
}
