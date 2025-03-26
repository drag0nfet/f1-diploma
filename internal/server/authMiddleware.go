package server

import (
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

// AuthMiddleware позволяет доступ всем, добавляя информацию об авторизации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username, _, _, response := services.CheckAuthCookie(r)
		if response.Success {
			w.Header().Set("X-Username", username)
		} else {
			w.Header().Set("X-Username", "") // Пустой логин для неавторизованных
		}
		// Передаём запрос дальше, даже если пользователь не авторизован
		next.ServeHTTP(w, r)
	}
}

// StrictAuthMiddleware разрешает доступ только авторизованным пользователям
func StrictAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username, _, _, response := services.CheckAuthCookie(r)
		if !response.Success {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}
		w.Header().Set("X-Username", username)
		next.ServeHTTP(w, r)
	}
}
