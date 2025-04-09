package handlers

import (
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	// Проверяем заголовок X-Requested-With
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

	// Проверяем метод
	if r.Method != http.MethodGet {
		response := services.Response{Success: false, Message: "Метод не поддерживается"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username, id, rights, response := services.CheckAuthCookie(r)
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := struct {
		Success  bool   `json:"success"`
		Username string `json:"username,omitempty"`
		Id       int    `json:"id,omitempty"`
		Rights   int    `json:"rights"`
		Message  string `json:"message,omitempty"`
	}{
		Success:  response.Success,
		Username: username,
		Id:       id,
		Rights:   rights,
		Message:  response.Message,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}
