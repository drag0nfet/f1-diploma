package userPage

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func GetMessagesList(w http.ResponseWriter, r *http.Request) {
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

	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указаны message_ids"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idStrings := strings.Split(idsParam, ",")
	if len(idStrings) == 0 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Список message_ids пуст"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var messageIDs []int
	for _, idStr := range idStrings {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный message_id: " + idStr})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		messageIDs = append(messageIDs, id)
	}

	var messages []models.Message
	if err := database.DB.
		Table("\"Message\"").
		Where("\"Message\".message_id IN ?", messageIDs).
		Find(&messages).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки сообщений"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success  bool             `json:"success"`
		Messages []models.Message `json:"messages"`
		Message  string           `json:"error_msg,omitempty"`
	}{
		Success:  true,
		Messages: messages,
		Message:  "",
	})
}
