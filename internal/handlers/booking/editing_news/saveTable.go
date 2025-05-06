package editing_news

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func SaveTable(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		HallID      int    `json:"hall_id"`
		TableID     int    `json:"table_id"`
		TableNamee  int    `json:"table_name"`
		PriceStatus string `json:"price_status"`
		SpotCount   int    `json:"spot_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	var table models.Table
	if err := database.DB.
		Where("table_name = ? AND hall_id = ?", req.TableNamee, req.HallID).
		First(&table).Error; err == nil && table.TableID != req.TableID {
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Стол с таким номером уже существует в этом зале",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.TableID == -1 {
		// Создание нового стола
		newTable := models.Table{
			HallID:      req.HallID,
			TableNamee:  req.TableNamee,
			PriceStatus: req.PriceStatus,
			Seats:       req.SpotCount,
		}
		if err := database.DB.Create(&newTable).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не удалось создать стол:" + err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Создание новых спотов
		var spots []models.Spot
		for number := 1; number <= newTable.Seats; number++ {
			spots = append(spots, models.Spot{
				TableID:  newTable.TableID,
				SpotName: number,
			})
		}

		if err := database.DB.Create(&spots).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не удалось создать споты: " + err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var tables []models.Table
		if err := database.DB.
			Table(`public."Table"`).
			Where("hall_id = ?", req.HallID).
			Order("seats DESC, table_name ASC").
			Find(&tables).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обновления столов"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(struct {
			Success bool           `json:"success"`
			Tables  []models.Table `json:"tables"`
		}{
			Success: true,
			Tables:  tables,
		})
		return
	}

	// Обновление существующего стола
	if err := database.DB.First(&table, req.TableID).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Стол не найден"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deltaSeats := req.SpotCount - table.Seats

	if deltaSeats > 0 {
		// Находим список всех существующих мест для данного стола
		var spotNames []int
		if err := database.DB.
			Model(&models.Spot{}).
			Where("table_id = ?", table.TableID).
			Pluck("spot_name", &spotNames).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не найдены места у данного стола"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Создаём мапу для быстрого доступа к существующим местам
		spotNamesMap := make(map[int]struct{})
		for _, spot := range spotNames {
			spotNamesMap[spot] = struct{}{}
		}

		// Создаём новые споты без совпадения имён
		var spots []models.Spot
		for number, deltaCount := 1, 0; deltaCount < deltaSeats; number++ {
			if _, found := spotNamesMap[number]; !found {
				deltaCount++
				spots = append(spots, models.Spot{
					TableID:  table.TableID,
					SpotName: number,
				})
			}
		}

		// Добавляем споты в БД
		if err := database.DB.Create(&spots).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не удалось создать споты: " + err.Error()})
			return
		}
	} else if deltaSeats < 0 {
		var spotIDs []uint
		if err := database.DB.
			Model(&models.Spot{}).
			Where("table_id = ?", table.TableID).
			Order("spot_name DESC").
			Limit(-deltaSeats).
			Pluck("spot_id", &spotIDs).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при поиске мест для удаления"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(spotIDs) > 0 {
			if err := database.DB.
				Where("spot_id IN ?", spotIDs).
				Delete(&models.Spot{}).Error; err != nil {
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при удалении мест"})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	table.PriceStatus = req.PriceStatus
	table.Seats = req.SpotCount
	table.TableNamee = req.TableNamee

	if err := database.DB.Save(&table).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не удалось обновить стол"})
		return
	}

	var tables []models.Table
	if err := database.DB.
		Table(`public."Table"`).
		Where("hall_id = ?", req.HallID).
		Order("seats DESC, table_name ASC").
		Find(&tables).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обновления столов"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Success bool           `json:"success"`
		Tables  []models.Table `json:"tables"`
	}{
		Success: true,
		Tables:  tables,
	})
}
