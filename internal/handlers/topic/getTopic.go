package topic

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetTopic(w http.ResponseWriter, r *http.Request) {
	// Проверяем заголовок X-Requested-With
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Проверяем метод
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем topicId из параметров маршрута
	vars := mux.Vars(r)
	topicIdStr, ok := vars["topicId"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "ID темы не указан"})
		return
	}

	topicId, err := strconv.Atoi(topicIdStr)
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный формат ID темы"})
		return
	}

	// Проверяем авторизацию
	_, _, _, response := services.CheckAuthCookie(r)
	if !response.Success {
		json.NewEncoder(w).Encode(response)
		return
	}

	// Ищем тему в базе
	var chat models.Chat
	if err := database.DB.Where("chat_id = ?", topicId).First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Тема не найдена"})
		} else {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки темы"})
		}
		return
	}

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success bool        `json:"success"`
		Topic   models.Chat `json:"topic"`
		Message string      `json:"message,omitempty"`
	}{
		Success: true,
		Topic:   chat,
		Message: "",
	})
}
