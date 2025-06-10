package topic

import (
	"diploma/internal/database"
	"diploma/internal/handlers"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"gopkg.in/guregu/null.v4"
	"net/http"
	"strconv"
	"time"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	// Проверяем заголовок X-Requested-With
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	// Проверяем метод
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	// Проверяем авторизацию, получаем userId для создания сообщения
	_, userId, rights, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	isConfirmed, err := handlers.GetConfirmation(userId)
	if err != nil {
		response = services.Response{Success: false, Message: "Ошибка при определении подтверждённости учётной записи"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isConfirmed {
		response = services.Response{Success: false, Message: "Учётная запись не подтверждена"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Нет бита 2
	if rights%4/2 == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Вы были заблокированы. Обратитесь к администратору",
		})
		return
	}

	// Декодируем тело запроса
	var req struct {
		ChatID  string `json:"chat_id"`
		Content string `json:"content"`
		ReplyID string `json:"reply_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный формат данных"})
		return
	}

	// Валидация данных
	chatID, err := strconv.Atoi(req.ChatID)
	if chatID <= 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при определении ID темы"})
		return
	}

	var replyID null.Int
	if req.ReplyID != "" {
		rid, err := strconv.Atoi(req.ReplyID)
		if err != nil || rid <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при определении reply_id"})
			return
		}

		// Проверяем, существует ли сообщение с таким reply_id
		var existingMessage models.Message
		if err := database.DB.Where("message_id = ?", rid).First(&existingMessage).Error; err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Сообщение, на которое вы отвечаете, не найдено"})
			return
		}

		replyID = null.IntFrom(int64(rid))
	}

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Текст сообщения не указан"})
		return
	}

	// Создаём новое сообщение
	preMessage := models.Message{
		ChatID:      chatID,
		Value:       req.Content,
		MessageTime: time.Now(),
		SenderID:    userId,
		ReplyID:     replyID,
	}

	// Сохраняем сообщение в базе
	if err := database.DB.Create(&preMessage).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Ошибка при сохранении сообщения",
		})
		return
	}

	if preMessage.MessageID%700000 == 650000 {
		services.NewPartition(preMessage.MessageID)
	}

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type MessageWithUsername struct {
		models.Message
		Username string `gorm:"column:login" json:"username"`
	}

	username, _, _, _ := services.CheckAuthCookie(r)
	message := MessageWithUsername{
		Message:  preMessage,
		Username: username,
	}

	json.NewEncoder(w).Encode(struct {
		Success  bool                `json:"success"`
		Message  MessageWithUsername `json:"message"`
		ErrorMsg string              `json:"error-msg,omitempty"`
	}{
		Success: true,
		Message: message,
	})
}
