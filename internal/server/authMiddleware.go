package server

import (
	"diploma/internal/services"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, response := services.CheckAuthCookie(r)
		// Добавляем информацию об авторизации в заголовки
		if response.Success {
			w.Header().Set("X-Username", username)
		} else {
			w.Header().Set("X-Username", "") // Пустой логин для неавторизованных
		}
		// Передаём запрос дальше, даже если пользователь не авторизован
		next.ServeHTTP(w, r)
	}
}
