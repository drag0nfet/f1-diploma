package userPage

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"time"
)

type BookingResponse struct {
	BookingID int          `json:"booking_id"`
	Status    string       `json:"status"`
	Event     models.Event `json:"event"`
	Table     models.Table `json:"table"`
	Spot      models.Spot  `json:"spot"`
	Hall      models.Hall  `json:"hall"`
}

func GetBookingInfo(w http.ResponseWriter, r *http.Request) []BookingResponse {
	// Проверяем авторизацию, получаем userId
	_, userId, _, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return nil
	}

	// Запрос к BookingSpot
	var bookingSpots []models.BookingSpot
	err := database.DB.
		Where("user_id = ? AND status = ? AND start_time >= ?", userId, "ACTIVE", time.Now().Truncate(24*time.Hour)).
		Order("start_time ASC").
		Find(&bookingSpots).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Ошибка при получении броней",
		})
		return nil
	}

	// Преобразуем в BookingResponse
	var bookings []BookingResponse
	for _, bs := range bookingSpots {
		var event models.Event
		var spot models.Spot
		var table models.Table
		var hall models.Hall

		// Загружаем Event
		if bs.EventID != nil {
			if err := database.DB.Where("event_id = ?", *bs.EventID).First(&event).Error; err != nil {
				continue // Пропускаем, если не удалось загрузить Event
			}
		} else {
			continue // Пропускаем, если EventID отсутствует
		}

		// Загружаем Spot
		if err := database.DB.Where("spot_id = ?", bs.SpotID).First(&spot).Error; err != nil {
			continue // Пропускаем, если не удалось загрузить Spot
		}

		// Загружаем Table
		if err := database.DB.Where("table_id = ?", spot.TableID).First(&table).Error; err != nil {
			continue // Пропускаем, если не удалось загрузить Table
		}

		// Загружаем Hall
		if err := database.DB.Where("hall_id = ?", table.HallID).First(&hall).Error; err != nil {
			continue // Пропускаем, если не удалось загрузить Hall
		}

		booking := BookingResponse{
			BookingID: bs.BookingID,
			Status:    bs.Status,
			Event:     event,
			Table:     table,
			Spot:      spot,
			Hall:      hall,
		}
		bookings = append(bookings, booking)
	}

	return bookings
}
