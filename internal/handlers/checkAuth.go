package handlers

import (
	"diploma/internal/database"
	"diploma/internal/services"
	"encoding/json"
	"errors"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		response := services.Response{
			Success: false,
			Message: "Прямой доступ к этому маршруту запрещён",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		response := services.Response{Success: false, Message: "Метод не поддерживается"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username, id, rights, response := services.CheckAuthCookie(r)

	isConfirmed, err := GetConfirmation(id)
	if err != nil {
		response = services.Response{Success: false, Message: "Ошибка при определении подтверждённости учётной записи"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := struct {
		Success     bool   `json:"success"`
		Username    string `json:"username,omitempty"`
		Id          int    `json:"id,omitempty"`
		Rights      int    `json:"rights"`
		Message     string `json:"message,omitempty"`
		IsConfirmed bool   `json:"is_confirmed"`
	}{
		Success:     response.Success,
		Username:    username,
		Id:          id,
		Rights:      rights,
		Message:     response.Message,
		IsConfirmed: isConfirmed,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}

func GetConfirmation(id int) (bool, error) {
	var isConfirmed bool
	err := database.DB.Table(`public."User"`).
		Select("is_confirmed").
		Where("user_id = ?", id).
		Scan(&isConfirmed).Error

	if err != nil {
		return false, errors.New("ошибка получения статуса подтверждения пользователя")
	}

	return isConfirmed, nil
}
