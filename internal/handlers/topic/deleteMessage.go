package topic

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	_, _, rights, response := services.CheckAuthCookie(r)
	if !response.Success {
		json.NewEncoder(w).Encode(response)
		return
	}

	if rights%2147483648 != 1 && rights%2 != 1 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Недостаточно прав"})
		return
	}

	vars := mux.Vars(r)
	messageIdStr, ok := vars["messageId"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "ID сообщения не указан"})
		return
	}
	messageID, err := strconv.ParseInt(messageIdStr, 10, 64)
	if err != nil || messageID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный ID сообщения"})
		return
	}

	var message models.Message
	if err := database.DB.Where("message_id = ?", messageID).First(&message).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Сообщение не найдено"})
		return
	}

	message.IsDeleted = true

	if err := database.DB.Save(&message).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Ошибка при сохранении удалённого сообщения",
		})
		return
	}
	/*
		// Старый вариант удаления, заменён в пользу скрытия
		if err := database.DB.Delete(&message).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(services.Response{
				Success: false,
				Message: "Ошибка при удалении сообщения",
			})
			return
		}
	*/

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Сообщение успешно удалено"})
}
