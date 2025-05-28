package booking

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func BookSpot(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим тело запроса
	var request struct {
		EventID int    `json:"event_id"`
		TableID int    `json:"table_id"`
		SpotID  *int   `json:"spot_id"`
		Action  string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(request, err)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат данных"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if request.EventID == 0 || request.TableID == 0 || request.Action == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Отсутствуют обязательные параметры"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if request.Action != "book" && request.Action != "cancel" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Недопустимое действие"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if request.Action == "book" && request.SpotID == nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указан spot_id для бронирования"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, userID, _, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Проверяем существование места и его принадлежность к столу
	var spot models.Spot
	if request.SpotID != nil {
		if err := database.DB.Table(`public."Spot"`).
			Where("spot_id = ? AND table_id = ?", *request.SpotID, request.TableID).
			First(&spot).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Место не найдено или не принадлежит указанному столу"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Проверяем существование события
	var event models.Event
	if err := database.DB.Table(`public."Event"`).
		Where("event_id = ?", request.EventID).
		First(&event).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Событие не найдено"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Обрабатываем действие
	if request.Action == "book" {
		// Проверяем, не забронировано ли место другим пользователем
		var existingBooking models.BookingSpot
		if err := database.DB.Table(`public."BookingSpot"`).
			Where("spot_id = ? AND event_id = ? AND status = ?", *request.SpotID, request.EventID, "ACTIVE").
			First(&existingBooking).Error; err == nil {
			if existingBooking.UserID != nil && *existingBooking.UserID != userID {
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Место уже забронировано другим пользователем"})
				w.WriteHeader(http.StatusConflict)
				return
			}
		}

		// Создаём или обновляем бронирование
		booking := models.BookingSpot{
			SpotID:    *request.SpotID,
			UserID:    &userID,
			EventID:   &request.EventID,
			Status:    "ACTIVE",
			StartTime: time.Now(),
		}
		if err := database.DB.Table(`public."BookingSpot"`).
			Where("spot_id = ? AND event_id = ?", *request.SpotID, request.EventID).
			Assign(booking).
			FirstOrCreate(&booking).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при создании бронирования"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if request.Action == "cancel" {
		// Отменяем бронирование только для текущего пользователя
		if err := database.DB.Table(`public."BookingSpot"`).
			Where("spot_id = ? AND event_id = ? AND user_id = ? AND status = ?", *request.SpotID, request.EventID, userID, "ACTIVE").
			Update("status", "INACTIVE").Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при отмене бронирования"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Формируем успешный ответ
	var mes string
	if request.Action == "book" {
		mes = "Место успешно забронировано"
	} else {
		mes = "Бронирование успешно отменено"
	}
	response = services.Response{
		Success: true,
		Message: mes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
