package topic

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"log"
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
	_, userId, _, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Декодируем тело запроса
	var req struct {
		ChatID  string `json:"chat_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(req, err.Error())
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

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Текст сообщения не указан"})
		return
	}

	// Создаём новое сообщение
	message := models.Message{
		ChatID:      chatID,
		Value:       req.Content,
		SenderID:    userId,
		MessageTime: time.Now(),
	}

	// Сохраняем сообщение в базе
	if err := database.DB.Create(&message).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Ошибка при сохранении сообщения",
		})
		return
	}

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(struct {
		Success  bool           `json:"success"`
		Message  models.Message `json:"message"`
		ErrorMsg string         `json:"error-msg,omitempty"`
	}{
		Success: true,
		Message: message,
	})
}
