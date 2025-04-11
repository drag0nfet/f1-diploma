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

func BlockUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	vars := mux.Vars(r)
	messageIdText, ok := vars["messageId"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Сообщение не указано"})
		return
	}

	messageId, err := strconv.Atoi(messageIdText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный message_id"})
		return
	}

	var user models.User

	err = database.DB.
		Table("Message").
		Select("\"User\".*").
		Joins("JOIN \"User\" ON \"Message\".sender_id = \"User\".user_id").
		Where("\"Message\".message_id = ?", messageId).
		Scan(&user).Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Пользователь не найден"})
		return
	}

	_, moderatorId, moderatorRights, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	if user.UserID == moderatorId {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Нельзя заблокировать самого себя"})
		return
	}

	if user.Rights%2 == 1 && moderatorRights/2147483648 != 1 {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Нельзя заблокировать модератора форума. Обратитесь к администратору",
		})
		return
	}

	user.Rights = user.Rights - user.Rights%4/2*2

	if err := database.DB.Save(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при обновлении прав пользователя"})
		return
	}

	blockNote := models.ForumBlockList{
		UserID:      user.UserID,
		MessageID:   messageId,
		ModeratorID: moderatorId,
	}

	if err := database.DB.Create(&blockNote).Error; err != nil {
		response := services.Response{Success: false, Message: "Ошибка при добавлении блокировки в список блокировок"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Пользователь успешно заблокирован"})
}
