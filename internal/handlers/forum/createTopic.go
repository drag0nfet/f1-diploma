package forum

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func CreateTopic(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный запрос"})
		return
	}

	// проверка авторизации - на практике всегда проходит успешно
	if _, _, _, response := services.CheckAuthCookie(r); !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	topic := models.Chat{
		Title:    &req.Title,
		ChatType: "forum",
		ItemID:   nil,
	}
	if err := database.DB.Create(&topic).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка создания темы"})
		return
	}

	json.NewEncoder(w).Encode(struct {
		Success bool   `json:"success"`
		TopicId int    `json:"topicId"`
		Message string `json:"message,omitempty"`
	}{Success: true, TopicId: topic.ChatID, Message: "Тема создана"})
}
