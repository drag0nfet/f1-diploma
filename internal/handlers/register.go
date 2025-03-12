package handlers

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := Response{Success: false, Message: "Метод не поддерживается"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{Success: false, Message: "Неверный формат данных"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var existingUser models.User
	if err := database.DB.Where("login = ?", req.Username).First(&existingUser).Error; err == nil {
		response := Response{Success: false, Message: "Пользователь с таким логином уже существует!"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusConflict)
		return
	}

	passHash, err := services.GetHash(req.Password)
	if err != nil {
		response := Response{Success: false, Message: "Ошибка при создании пароля"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := models.User{Login: req.Username, Password: passHash}
	if err := database.DB.Create(&user).Error; err != nil {
		response := Response{Success: false, Message: "Ошибка при регистрации"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Успешная регистрация пользователя", req.Username)
	response := Response{Success: true}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
