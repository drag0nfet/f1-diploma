package userPage

import (
	"diploma/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetBookingPass(w http.ResponseWriter, r *http.Request) {
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

	_, userId, _, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	ans := createPassQR(bookings, userId)

	// Формируем ответ
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(ans)
}

func createPassQR(bookings []BookingResponse, userId int) []byte {
	if len(bookings) == 0 {
		return nil
	}

	// Собираем все spot_id в слайс строк
	spotIDs := make([]string, 0, len(bookings))
	for _, booking := range bookings {
		spotIDs = append(spotIDs, strconv.Itoa(booking.Spot.SpotID))
	}

	code := fmt.Sprintf("U%dS%s", userId, strings.Join(spotIDs, ","))
	codeHash, err := services.GetHash(code)
	if err != nil {
		return nil
	}
	url := "https://api.qrserver.com/v1/create-qr-code/?data=" + codeHash + "&size=200x200"

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	qrImage, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return qrImage
}
