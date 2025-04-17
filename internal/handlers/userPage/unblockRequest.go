package userPage

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func UnblockRequest(w http.ResponseWriter, r *http.Request) {
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
		UserID    int    `json:"user_id"`
		MessageID int    `json:"message_id"`
		Comment   string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный запрос"})
		return
	}

	unblockRequest := models.UnblockRequest{
		UserID:    req.UserID,
		MessageID: req.MessageID,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(&unblockRequest).Error; err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка создания запроса на разблокировку"})
		return
	}

	if err := database.DB.
		Table("\"ForumBlockList\"").
		Where("user_id = ? AND message_id = ?", req.UserID, req.MessageID).
		Update("status", "WAITING").Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обновления статуса блокировки"})
		return
	}

	json.NewEncoder(w).Encode(struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}{Success: true, Message: "Запрос на разблокировку отправлен"})
}
