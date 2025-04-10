package bar

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"encoding/json"
	"net/http"
)

func GetDishes(w http.ResponseWriter, r *http.Request) {
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

	var dishes []models.Dish
	if err := database.DB.Find(&dishes).Error; err != nil {
		json.NewEncoder(w).Encode(services.Response{Success: false, Message: "Ошибка загрузки блюд"})
		return
	}

	responseDishes := make([]struct {
		DishID      int    `json:"dish_id"`
		Name        string `json:"name"`
		Cost        int    `json:"cost"`
		Description string `json:"description"`
		Image       []byte `json:"image"`
	}, len(dishes))

	for i, dish := range dishes {
		responseDishes[i] = struct {
			DishID      int    `json:"dish_id"`
			Name        string `json:"name"`
			Cost        int    `json:"cost"`
			Description string `json:"description"`
			Image       []byte `json:"image"`
		}{
			DishID:      dish.DishID,
			Name:        dish.Name,
			Cost:        dish.Cost,
			Description: dish.Description,
			Image:       dish.Image,
		}
	}

	json.NewEncoder(w).Encode(struct {
		Success bool `json:"success"`
		Dishes  []struct {
			DishID      int    `json:"dish_id"`
			Name        string `json:"name"`
			Cost        int    `json:"cost"`
			Description string `json:"description"`
			Image       []byte `json:"image"`
		} `json:"dishes"`
	}{Success: true, Dishes: responseDishes})
}
