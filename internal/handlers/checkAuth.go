package handlers

import (
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	username, response := services.CheckAuthCookie(r)
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := struct {
		Success  bool   `json:"success"`
		Username string `json:"username,omitempty"`
		Message  string `json:"message,omitempty"`
	}{
		Success:  response.Success,
		Username: username,
		Message:  response.Message,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}
