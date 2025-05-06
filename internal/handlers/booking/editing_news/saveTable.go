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
		return
	}

	if req.TableID == -1 {
		// Создание новой записи
		newTable := models.Table{
			HallID:      req.HallID,
			TableNamee:  req.TableNamee,
			PriceStatus: req.PriceStatus,
			Seats:       req.SpotCount,
		}
		if err := database.DB.Create(&newTable).Error; err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не удалось создать стол:" + err.Error()})
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
		return
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
