package account

import (
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// Удаляем куки, устанавливая её срок действия в прошлое
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

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	response := services.Response{
		Success: true,
		Message: "Вы успешно вышли из системы",
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
