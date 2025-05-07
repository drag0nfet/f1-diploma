package editing_halls

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func SaveHall(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Прямой доступ запрещён"})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг формы с ограничением 32 МБ
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка парсинга формы"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hallIdStr := r.FormValue("hall_id")
	name := r.FormValue("name")
	description := r.FormValue("description")
	deletedRaw := r.FormValue("deleted_photo_ids") // Изменено с deleted_indices

	// Валидация hall_id
	hallId, err := strconv.Atoi(hallIdStr)
	if err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный hall_id"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if name == "" || description == "" {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Название и описание обязательны"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Обработка удалённых фотографий
	var deletedPhotoIds []int
	if deletedRaw != "" {
		if err = json.Unmarshal([]byte(deletedRaw), &deletedPhotoIds); err != nil {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Некорректный формат deleted_photo_ids"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Обработка новых фотографий
	form := r.MultipartForm
	files := form.File["photos"]
	var newPhotos []models.HallPhoto

	for _, fileHeader := range files {
		// Валидация размера файла (макс. 5 МБ)
		if fileHeader.Size > 5<<20 {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Файл слишком большой (макс. 5 МБ)"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Проверка MIME-типа
		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			continue
		}
		mimeType := http.DetectContentType(buffer[:n])
		if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/bmp" && mimeType != "image/jpg" && mimeType != "image/webp" {
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Поддерживаются только JPEG, JPG, BMP, WEBP и PNG"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Сброс чтения файла
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			log.Printf("Ошибка сброса чтения файла %s: %v", fileHeader.Filename, err)
			continue
		}

		content, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Ошибка чтения файла %s: %v", fileHeader.Filename, err)
			continue
		}

		newPhotos = append(newPhotos, models.HallPhoto{
			HallID:   hallId,
			Content:  content,
			MimeType: mimeType,
		})
	}

	// Создание нового зала
	if hallId == -1 {
		newHall := models.Hall{
			Name:        name,
			Description: description,
		}
		if err := database.DB.Create(&newHall).Error; err != nil {
			log.Printf("Ошибка создания зала: %v", err)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка создания зала"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Сохранение фотографий
		for i := range newPhotos {
			newPhotos[i].HallID = newHall.HallID
			if err := database.DB.Create(&newPhotos[i]).Error; err != nil {
				log.Printf("Ошибка сохранения фотографии для зала %d: %v", newHall.HallID, err)
				json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка сохранения фотографий"})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// Загрузка фотографий для ответа
		var photos []models.HallPhoto
		if err := database.DB.
			Table(`public."HallPhotos"`).
			Where("hall_id = ?", newHall.HallID).
			Order("created_at ASC").
			Find(&photos).Error; err != nil {
		}

		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"hall": map[string]any{
				"hall_id":     newHall.HallID,
				"name":        newHall.Name,
				"description": newHall.Description,
				"album":       photos,
			},
		})
		return
	}

	// Обновление существующего зала
	var hall models.Hall
	if err := database.DB.
		Table(`public."Hall"`).
		First(&hall, hallId).Error; err != nil {
		log.Printf("Зал %d не найден: %v", hallId, err)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Зал не найден"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Удаление фотографий
	if len(deletedPhotoIds) > 0 {
		if err := database.DB.
			Table(`public."HallPhotos"`).
			Where("hall_id = ? AND id IN ?", hallId, deletedPhotoIds).
			Delete(&models.HallPhoto{}).Error; err != nil {
			log.Printf("Ошибка удаления фотографий для зала %d: %v", hallId, err)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка удаления фотографий"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Обновление данных зала
	hall.Name = name
	hall.Description = description
	if err := database.DB.Save(&hall).Error; err != nil {
		log.Printf("Ошибка обновления зала %d: %v", hallId, err)
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка обновления зала"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Сохранение новых фотографий
	for i := range newPhotos {
		if err := database.DB.Create(&newPhotos[i]).Error; err != nil {
			log.Printf("Ошибка сохранения фотографии для зала %d: %v", hallId, err)
			json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка сохранения фотографий"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Загрузка актуальных фотографий для ответа
	var photos []models.HallPhoto
	if err := database.DB.
		Table(`public."HallPhotos"`).
		Where("hall_id = ?", hallId).
		Order("created_at ASC").
		Find(&photos).Error; err != nil {
		log.Printf("Ошибка загрузки фотографий для зала %d: %v", hallId, err)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"hall": map[string]any{
			"hall_id":     hall.HallID,
			"name":        hall.Name,
			"description": hall.Description,
			"album":       photos,
		},
	})
}
