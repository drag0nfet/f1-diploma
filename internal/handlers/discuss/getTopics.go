package discuss

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"time"
)

func GetTopics(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var topics []models.Chat
	if err := database.DB.Find(&topics).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки тем"})
		return
	}

	responseTopics := make([]struct {
		ChatID    int       `json:"chat_id"`
		ChatType  string    `json:"chat_type"`
		ItemID    *int      `json:"item_id"`
		CreatedAt time.Time `json:"created_at"`
		Title     string    `json:"title"`
	}, len(topics))

	for i, topic := range topics {
		title := ""
		if topic.Title != nil {
			title = *topic.Title
		}
		responseTopics[i] = struct {
			ChatID    int       `json:"chat_id"`
			ChatType  string    `json:"chat_type"`
			ItemID    *int      `json:"item_id"`
			CreatedAt time.Time `json:"created_at"`
			Title     string    `json:"title"`
		}{
			ChatID:    topic.ChatID,
			ChatType:  topic.ChatType,
			ItemID:    topic.ItemID,
			CreatedAt: topic.CreatedAt,
			Title:     title,
		}
	}

	json.NewEncoder(w).Encode(struct {
		Success bool `json:"success"`
		Topics  []struct {
			ChatID    int       `json:"chat_id"`
			ChatType  string    `json:"chat_type"`
			ItemID    *int      `json:"item_id"`
			CreatedAt time.Time `json:"created_at"`
			Title     string    `json:"title"`
		} `json:"topics"`
	}{Success: true, Topics: responseTopics})
}
