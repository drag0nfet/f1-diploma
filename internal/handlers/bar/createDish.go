package bar

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

func CreateDish(w http.ResponseWriter, r *http.Request) {
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

	_, _, rights, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	if rights%8/4 != 1 {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Только модераторы могут добавлять блюда"})
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Максимум 10MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обработки формы"})
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	costStr := r.FormValue("cost")

	file, _, err := r.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки файла"})
		return
	}

	if name == "" || costStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Название и цена обязательны!"})
		return
	}
	cost, err := strconv.Atoi(costStr)
	if err != nil || cost < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Неверная цена"})
		return
	}

	var imageData []byte
	if file != nil {
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при обработке файла"})
				return
			}
		}(file)
		imageData, err = io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка чтения файла"})
			return
		}
	}

	dish := models.Dish{
		Name:        name,
		Cost:        cost,
		Description: description,
		Image:       imageData,
	}

	if err := database.DB.Create(&dish).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка сохранения блюда"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.Response{Success: true, Message: "Блюдо успешно добавлено"})
}
