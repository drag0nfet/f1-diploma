package index

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := services.Response{Success: false, Message: "Метод не поддерживается"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var auth LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		response := services.Response{Success: false, Message: "Неверный формат данных"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.Where("login = ?", auth.Username).First(&user).Error; err != nil {
		response := services.Response{Success: false, Message: "Пользователь с таким логином не найден!"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !services.CheckHash(auth.Password, user.Password) {
		response := services.Response{Success: false, Message: "Неверный пароль!"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie := services.NewCookie(w, auth.Username, user.Rights, user.UserID)
	if cookie.Name == "" {
		return // Ошибка уже обработана внутри NewCookie
	}
	http.SetCookie(w, &cookie)

	response := services.Response{Success: true}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
