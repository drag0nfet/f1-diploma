package topic

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func BlockUser(w http.ResponseWriter, r *http.Request) {
	log.Println("BlockUser")
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
	username, ok := vars["username"]
	if !ok {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Пользователь не указан"})
		return
	}

	var user models.User
	if err := database.DB.Where("login = ?", username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Пользователь не найден"})
		return
	}

	if _, _, _, response := services.CheckAuthCookie(r); !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, moderatorId, _, _ := services.CheckAuthCookie(r)
	if user.UserID == moderatorId {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Нельзя заблокировать самого себя"})
		return
	}

	user.Rights = user.Rights - user.Rights%4/2*2

	if err := database.DB.Save(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при обновлении прав пользователя"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Пользователь успешно заблокирован"})
}
