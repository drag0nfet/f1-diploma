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
	if err := database.DB.Where("login = ? OR email = ?", auth.Username, auth.Username).
		First(&user).Error; err != nil {
		response := services.Response{Success: false, Message: "Пользователь с таким логином или email не найден!"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !services.CheckHash(auth.Password, user.Password) {
		response := services.Response{Success: false, Message: "Неверный пароль!"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	/*
		isConfirmed := user.IsConfirmed
		if !isConfirmed {
			response := services.Response{Success: false, Message: "Учётная запись не подтверждена!"}
			json.NewEncoder(w).Encode(response)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	*/

	// Считаем блокировки на форуме
	var blockCount int64
	if err := database.DB.
		Table("\"ForumBlockList\"").
		Where("user_id = ? AND is_valid = ?", user.UserID, true).
		Count(&blockCount).Error; err != nil {
		response := services.Response{Success: false, Message: "Ошибка проверки блокировок"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Если нет активных блокировок (blockCount == 0), устанавливаем второй бит в 1
	if blockCount == 0 {
		if user.Rights%4 < 2 { // Если второй бит 0 (rights%4 = 0 или 1)
			user.Rights += 2 // Устанавливаем второй бит в 1 (rights%4 станет 2 или 3)
		}
	} else {
		if user.Rights%4 >= 2 { // Если второй бит 1 (rights%4 = 2 или 3)
			user.Rights -= 2 // Устанавливаем второй бит в 0 (rights%4 станет 0 или 1)
		}
	}

	// Обновляем права пользователя в базе
	if err := database.DB.Save(&user).Error; err != nil {
		response := services.Response{Success: false, Message: "Ошибка обновления прав пользователя"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
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
