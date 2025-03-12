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

func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("reg try")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Метод не поддерживается"})
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Неверный формат данных"})
		return
	}

	log.Printf("Получены регистрационные данные на логин %s\n", req.Username)

	passHash, err := services.GetHash(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка при создании пароля"})
		return
	}

	user := models.User{Login: req.Username, Password: passHash}
	if err := database.DB.Create(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка при регистрации"})
		return
	}

	log.Println("Успешная регистрация пользователя", req.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
