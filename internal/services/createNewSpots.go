package services

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"fmt"
	"time"
)

func CreateOrUpdateSpots(eventId int, timeStart time.Time, isNew bool) error {
	var spotIds []int
	if err := database.DB.Table(`public."Spot"`).Pluck("spot_id", &spotIds).Error; err != nil {
		return fmt.Errorf("ошибка нахождения спотов: %v", err)
	}

	if len(spotIds) == 0 {
		return fmt.Errorf("нет доступных спотов")
	}

	if isNew {
		var bookingSpots []models.BookingSpot
		for _, spotID := range spotIds {
			bookingSpots = append(bookingSpots, models.BookingSpot{
				SpotID:    spotID,
				EventID:   &eventId,
				StartTime: timeStart,
				// UserID по умолчанию nil
				// Status по умолчанию INACTIVE (установлено на уровне БД)
			})
		}
		if err := database.DB.Table(`public."BookingSpot"`).Create(&bookingSpots).Error; err != nil {
			return fmt.Errorf("ошибка создания букинг-спотов: %v", err)
		}
	} else {
		if err := database.DB.Table(`public."BookingSpot"`).
			Where("event_id = ?", eventId).
			Updates(map[string]interface{}{"start_time": timeStart}).Error; err != nil {
			return fmt.Errorf("ошибка обновления времени букинг-спотов: %v", err)
		}
	}

	return nil
}
