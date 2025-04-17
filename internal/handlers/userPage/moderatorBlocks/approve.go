package moderatorBlocks

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func Approve(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		return
	}

	vars := mux.Vars(r)
	request_idText, ok := vars["request_id"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указан номер запроса"})
		return
	}

	request_id, err := strconv.Atoi(request_idText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный request_id"})
		return
	}

	var req struct {
		UserID    int `json:"user_id"`
		MessageID int `json:"message_id"`
	}
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверный формат данных"})
		return
	}

	var block models.ForumBlockList
	if err = database.DB.Where("message_id = ? AND user_id = ?", req.MessageID, req.UserID).First(&block).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Блокировка не найдена"})
		return
	}

	block.IsValid = false

	if err = database.DB.Save(&block).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Ошибка при сохранении снятой блокировки",
		})
		return
	}

	var request models.UnblockRequest
	if err = database.DB.Where("request_id = ?", request_id).First(&request).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Запрос на разблокировку не найден"})
		return
	}

	request.Status = "APPROVED"

	if err = database.DB.Save(&request).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{
			Success: false,
			Message: "Ошибка при сохранении принятой апелляции",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Успешное принятие апелляции"})
}
