package userPage

import (
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func Logout(w http.ResponseWriter, _ *http.Request) {
	cookie := http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Удаляем куки
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	response := services.Response{
		Success: true,
		Message: "Вы успешно вышли из системы",
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
