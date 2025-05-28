package userPage

import (
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func GetBookings(w http.ResponseWriter, r *http.Request) {
	// Проверяем заголовок X-Requested-With
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	// Проверяем метод
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	bookings := GetBookingInfo(w, r)
	if bookings == nil {
		return
	}

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Success bool              `json:"success"`
		Data    []BookingResponse `json:"data"`
	}{
		Success: true,
		Data:    bookings,
	})
}
