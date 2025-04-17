package moderatorBlocks

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func GetRequests(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Метод не поддерживается"})
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requests []models.UnblockRequest
	err := database.DB.Find(&requests).Error
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки запросов на разблокировку"})
		return
	}

	json.NewEncoder(w).Encode(struct {
		Success  bool                    `json:"success"`
		Requests []models.UnblockRequest `json:"requests"`
	}{Success: true, Requests: requests})
}
