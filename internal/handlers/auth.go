package handlers

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Попытка логина")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Метод не поддерживается"})
		return
	}

	var auth LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Неверный формат данных"})
		return
	}

	var user models.User
	if err := database.DB.Where("login = ?", auth.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь с таким логином не найден!"})
		return
	}

	if !services.CheckHash(auth.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Неверный пароль!"})
		return
	}

	// Можно установить cookie или токен
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
