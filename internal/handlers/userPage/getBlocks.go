package userPage

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GetBlocks(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Не указан username"})
		return
	}

	var blocks []models.ForumBlockList
	err := database.DB.
		Joins("JOIN \"User\" ON \"ForumBlockList\".user_id = \"User\".user_id").
		Where("\"User\".login = ?", username).
		Order("time_got DESC").
		Find(&blocks).Error
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки блокировок"})
		return
	}

	json.NewEncoder(w).Encode(struct {
		Success bool                    `json:"success"`
		Blocks  []models.ForumBlockList `json:"blocks"`
	}{Success: true, Blocks: blocks})
}
