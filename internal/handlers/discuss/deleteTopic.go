package discuss

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func DeleteDiscuss(w http.ResponseWriter, r *http.Request) {
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

	// Проверяем авторизацию и права
	_, response, rights := services.CheckAuthCookie(r)
	if !response.Success {
		json.NewEncoder(w).Encode(response)
		return
	}

	if rights%2 != 1 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Недостаточно прав"})
		return
	}

	// Декодируем тело запроса
	var req struct {
		ChatID int `json:"chat_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный запрос"})
		return
	}

	// Удаляем тему из базы
	if err := database.DB.Delete(&models.Chat{}, req.ChatID).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка удаления темы"})
		return
	}

	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Тема удалена"})
}
