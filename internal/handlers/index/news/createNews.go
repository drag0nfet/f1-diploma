package news

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"time"
)

func CreateNews(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_, UserID, _, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Comment     string  `json:"comment"`
		Image       []byte  `json:"image"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный запрос"})
		return
	}

	news := models.News{
		Title:       req.Title,
		CreatorID:   UserID,
		Description: req.Description,
		Comment:     req.Comment,
		CreatedAt:   time.Now(),
		Image:       req.Image,
	}
	if err := database.DB.Create(&news).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка создания новости"})
		return
	}

	json.NewEncoder(w).Encode(struct {
		Success bool        `json:"success"`
		News    models.News `json:"editing_news"`
		Message string      `json:"message,omitempty"`
	}{Success: true, News: news, Message: "Новость создана"})
}
