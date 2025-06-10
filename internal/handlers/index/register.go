package index

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"time"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := services.Response{Success: false, Message: "Метод не поддерживается"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := services.Response{Success: false, Message: "Неверный формат данных"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var existingUser models.User
	if err := database.DB.Where("login = ?", req.Username).First(&existingUser).Error; err == nil {
		response := services.Response{Success: false, Message: "Пользователь с таким логином уже существует!"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		response := services.Response{Success: false, Message: "Пользователь с такой почтой уже существует!"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusConflict)
		return
	}

	passHash, err := services.GetHash(req.Password)
	if err != nil {
		response := services.Response{Success: false, Message: "Ошибка при хешировании пароля"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	confToken, err := services.SendConfirmation(req.Email)
	if err != nil {
		response := services.Response{Success: false, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := models.User{
		Login:             req.Username,
		Password:          passHash,
		Email:             req.Email,
		ConfirmationToken: confToken,
		LastSent:          time.Now(),
	}
	if err = database.DB.Create(&user).Error; err != nil {
		response := services.Response{Success: false, Message: "Ошибка при регистрации"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := services.Response{Success: true}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
