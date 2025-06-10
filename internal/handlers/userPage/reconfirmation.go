package userPage

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"time"
)

func Reconfirmation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := services.Response{Success: false, Message: "Метод не поддерживается"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")

	var user models.User

	if err := database.DB.Where("login = ?", username).First(&user).Error; err != nil {
		response := services.Response{Success: false, Message: "Пользователь не найден"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if time.Since(user.LastSent) < 1*time.Hour {
		response := services.Response{
			Success: false,
			Message: "Получить новую ссылку для подтверждения можно через час после получения последнего письма"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := services.SendConfirmation(user.Email)
	if err != nil {
		response := services.Response{Success: false, Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err = database.DB.Model(&user).Updates(models.User{
		ConfirmationToken: token,
		LastSent:          time.Now(),
	}).Error; err != nil {
		response := services.Response{Success: false, Message: "Ошибка подтверждения"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := services.Response{Success: true, Message: "Отправили письмо для подтверждения"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
