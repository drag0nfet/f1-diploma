package index

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
	"time"
)

func Confirm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	var user models.User

	if err := database.DB.Where("confirmation_token = ?", token).First(&user).Error; err != nil {
		response := services.Response{Success: false, Message: "Недействительный токен"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if time.Since(user.LastSent) > 1*time.Hour {
		response := services.Response{Success: false, Message: "Срок действия токена истёк. Запросите новый в ЛК"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DB.Model(&user).Updates(models.User{
		IsConfirmed:       true,
		ConfirmationToken: "",
	}).Error; err != nil {
		response := services.Response{Success: false, Message: "Ошибка подтверждения"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	cookie := services.NewCookie(w, user.Login, user.Rights, user.UserID)
	if cookie.Name == "" {
		return // Ошибка уже обработана внутри NewCookie
	}
	http.SetCookie(w, &cookie)

	http.ServeFile(w, r, "dir/pages/index.html")
}
