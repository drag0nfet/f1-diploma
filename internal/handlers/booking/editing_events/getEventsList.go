package editing_events

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func GetEventsList(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.Error(w, "Прямой доступ запрещён", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем параметры запроса
	dateFrom := r.URL.Query().Get("date_from")   // Формат дд.мм.гггг
	dateRange := r.URL.Query().Get("date_range") // Формат дд.мм.гггг-дд.мм.гггг
	sportCategory := r.URL.Query().Get("sport_category")

	currentDate := time.Now()

	// Инициализируем запрос к базе
	var events interface{} // Используем interface{}, чтобы поддерживать разные структуры результата
	var err error

	if dateFrom == "" && dateRange == "" && sportCategory == "" {
		// Простой список ивентов (только event_id и description)
		var simpleEvents []struct {
			EventID     int    `json:"event_id"`
			Description string `json:"description"`
		}
		err = database.DB.
			Table(`public."Event"`).
			Select("event_id, description").
			Order("event_id ASC").
			Find(&simpleEvents).Error
		events = simpleEvents
	} else {
		// Расширенный список с фильтрами
		var filteredEvents []models.Event
		query := database.DB.Table(`public."Event"`)

		// Обработка фильтров по дате
		if dateRange != "" {
			// Парсим диапазон дат
			dates := strings.Split(dateRange, "—")
			if len(dates) != 2 {
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат диапазона дат"})
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			startDateStr := strings.TrimSpace(dates[0])
			endDateStr := strings.TrimSpace(dates[1])

			const dateFormat = "02.01.2006"
			startDate, err := time.Parse(dateFormat, startDateStr)
			if err != nil {
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат начальной даты"})
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			endDate, err := time.Parse(dateFormat, endDateStr)
			if err != nil {
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат конечной даты"})
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Устанавливаем время на начало дня для начальной даты и конец дня для конечной
			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)

			query = query.Where("time_start BETWEEN ? AND ?", startDate, endDate)
		} else if dateFrom != "" {
			// Парсим дату начала в формате дд.мм.гггг
			const dateFormat = "02.01.2006"
			startDate, err := time.Parse(dateFormat, dateFrom)
			if err != nil {
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат начальной даты"})
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			query = query.Where("time_start >= ?", startDate)
		} else {
			// Дефолт: с текущей даты и позже
			query = query.Where("time_start >= ?", currentDate)
		}

		// Обработка фильтра по категории спорта
		if sportCategory != "" {
			categories := strings.Split(sportCategory, ",")
			if len(categories) > 0 {
				query = query.Where("sport_category IN ?", categories)
			}
		}

		// Сортировка по time_start
		query = query.Order("time_start ASC")

		// Выполняем запрос
		err = query.Find(&filteredEvents).Error
		events = filteredEvents
	}

	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки событий"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := map[string]any{
		"success": true,
		"events":  events,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
