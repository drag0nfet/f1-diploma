package editing_events

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func SaveEvent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат данных"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Description == "" || event.SportCategory == "" || event.SportType == "" || event.PriceStatus == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Все поля обязательны для заполнения"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Duration <= 0 {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Продолжительность должна быть больше 0"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isNew := event.EventID == -1

	if !isNew {
		var count int64
		if err := database.DB.Table(`public."Event"`).Where("event_id = ?", event.EventID).Count(&count).Error; err != nil || count == 0 {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ивент не найден"})
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	if isNew {
		var existingEvents []models.Event
		endTime := event.TimeStart.Add(time.Duration(event.Duration) * time.Minute)
		startWindow := event.TimeStart.AddDate(0, 0, -1)
		endWindow := event.TimeStart.AddDate(0, 0, 1)

		if err := database.DB.
			Table(`public."Event"`).
			Where("event_id != ?", event.EventID).
			Where("time_start BETWEEN ? AND ?", startWindow, endWindow).
			Find(&existingEvents).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка проверки расписания"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var conflicts []string
		for _, existing := range existingEvents {
			existingEnd := existing.TimeStart.Add(time.Duration(existing.Duration) * time.Minute)
			if (event.TimeStart.Before(existingEnd) && endTime.After(existing.TimeStart)) ||
				(event.TimeStart.Equal(existing.TimeStart) || endTime.Equal(existingEnd)) {
				conflictStr := fmt.Sprintf("%s до %s",
					existing.TimeStart.Format("15:04"),
					existingEnd.Format("15:04"))
				conflicts = append(conflicts, conflictStr)
			}
		}

		if len(conflicts) > 0 {
			message := "Время с " + strings.Join(conflicts, ", ") + " занято другими ивентами. Измените время в соответствии с расписанием!"
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: message})
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	if isNew {
		newEvent := models.Event{
			Description:   event.Description,
			TimeStart:     event.TimeStart,
			SportCategory: event.SportCategory,
			SportType:     event.SportType,
			PriceStatus:   event.PriceStatus,
			Duration:      event.Duration,
		}
		if err := database.DB.
			Table(`public."Event"`).
			Create(&newEvent).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка создания ивента"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		event.EventID = newEvent.EventID
	} else {
		if err := database.DB.
			Table(`public."Event"`).
			Where("event_id = ?", event.EventID).
			Updates(&event).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обновления ивента"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := services.CreateOrUpdateSpots(event.EventID, event.TimeStart, isNew); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"success":  true,
		"event_id": event.EventID,
		"message":  "Ивент успешно сохранён",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
