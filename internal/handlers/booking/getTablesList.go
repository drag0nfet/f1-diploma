package booking

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func GetTablesList(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем hallId из пути и event_id из параметров запроса
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат URL"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hallIdStr := pathParts[3]
	hallId, err := strconv.Atoi(hallIdStr)

	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный hall_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventIdStr := r.URL.Query().Get("event_id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный event_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получаем столы для указанного зала
	var tables []models.Table
	if err := database.DB.
		Table(`public."Table"`).
		Where("hall_id = ?", hallId).
		Find(&tables).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки столов"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type SpotResponse struct {
		SpotID   int                  `json:"spot_id"`
		TableID  int                  `json:"table_id"`
		SpotName string               `json:"spot_name"`
		Bookings []models.BookingSpot `json:"bookings"`
	}

	// Подготавливаем ответ с вложенными спотами и бронированиями
	type TableResponse struct {
		TableID     int            `json:"table_id"`
		HallID      int            `json:"hall_id"`
		TableName   string         `json:"table_name"`
		PriceStatus string         `json:"price_status"`
		Seats       int            `json:"seats"`
		Spots       []SpotResponse `json:"spots"`
	}

	var tableResponses []TableResponse
	for _, table := range tables {
		var spots []models.Spot
		if err = database.DB.
			Table(`public."Spot"`).
			Where("table_id = ?", table.TableID).
			Find(&spots).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки спотов"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Формируем споты с бронированиями
		var spotResponses []SpotResponse
		for _, spot := range spots {
			// Получаем бронирования для спота и ивента
			var bookings []models.BookingSpot
			if err = database.DB.
				Table(`public."BookingSpot"`).
				Where("spot_id = ? AND event_id = ?", spot.SpotID, eventId).
				Find(&bookings).Error; err != nil {
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки бронирований"})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			spotResponses = append(spotResponses, SpotResponse{
				SpotID:   spot.SpotID,
				TableID:  spot.TableID,
				SpotName: strconv.Itoa(spot.SpotName),
				Bookings: bookings,
			})
		}

		tableResponses = append(tableResponses, TableResponse{
			TableID:     table.TableID,
			HallID:      table.HallID,
			TableName:   strconv.Itoa(table.TableNamee),
			PriceStatus: table.PriceStatus,
			Seats:       table.Seats,
			Spots:       spotResponses,
		})
	}

	// Формируем итоговый ответ
	response := map[string]any{
		"success": true,
		"tables":  tableResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
