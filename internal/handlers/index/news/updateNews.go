package news

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

func UpdateNews(w http.ResponseWriter, r *http.Request) {
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

	_, userId, userRights, response := services.CheckAuthCookie(r)
	if !response.Success {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	isModerator := userRights%16 >= 8 || userRights/2147483648 == 1
	if !isModerator {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Доступ запрещён"})
		return
	}

	err := r.ParseMultipartForm(32 << 20) // Лимит 32MB
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обработки данных формы"})
		return
	}

	newsIdStr := r.FormValue("news_id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	comment := r.FormValue("comment")
	status := r.FormValue("status")

	newsId, err := strconv.Atoi(newsIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный news_id"})
		return
	}

	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Название новости обязательно"})
		return
	}

	if status != "ACTIVE" && status != "ARCHIVE" && status != "DRAFT" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный статус новости"})
		return
	}

	var news models.News
	if newsId >= 0 {
		// Обновление существующей новости
		if err = database.DB.Where("news_id = ?", newsId).First(&news).Error; err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Новость не найдена"})
			return
		}

		// Проверяем, что пользователь — создатель новости или администратор
		if news.CreatorID != userId && userRights/2147483648 != 1 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "У вас нет прав для редактирования этой новости"})
			return
		}
	} else {
		// Создание новой новости
		news = models.News{
			CreatorID: userId,
		}
	}

	// Обновляем поля новости
	news.Title = title
	news.Description = &description
	news.Comment = comment
	news.Status = status

	// Обработка изображения
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageData, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка чтения изображения"})
			return
		}
		news.Image = imageData
	} else if !errors.Is(err, http.ErrMissingFile) {
		// Если ошибка не связана с отсутствием файла, возвращаем ошибку
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки изображения"})
		return
	}

	// Сохранение новости
	if newsId >= 0 {
		// Обновление существующей новости
		if err := database.DB.Save(&news).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при обновлении новости"})
			return
		}
	} else {
		// Создание новой новости
		if err := database.DB.Create(&news).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка при создании новости"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Success bool   `json:"success"`
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		Success: true,
		ID:      news.NewsID,
		Message: "Новость успешно создана или изменена",
	})
}
