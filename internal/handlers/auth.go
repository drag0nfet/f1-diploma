package handlers

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"diploma/internal/services"
	"fmt"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		passHash, err := services.GetHash(password)

		if err != nil {
			fmt.Fprintf(w, "Ошибка при создании пароля: %v", err)
			return
		}

		user := models.User{Login: username, Password: passHash}

		result := database.DB.Create(&user)
		if result.Error != nil {
			fmt.Fprintf(w, "Ошибка при регистрации: %v", result.Error)
			return
		}
		log.Println("Успешная регистрация пользователя ", username)
		fmt.Fprintf(w, "Регистрация успешна!")
	} else {
		http.ServeFile(w, r, "web/index.html")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user models.User
		userCheck := database.DB.Where("login = ?", username).First(&user)
		if userCheck.Error != nil {
			fmt.Fprintf(w, "Пользователь с таким логином не найден!")
			return
		}

		if !services.CheckHash(password, user.Password) {
			fmt.Fprintf(w, "Неверный пароль!")
			return
		}

		fmt.Fprintf(w, "Авторизация успешна!")
	} else {
		http.ServeFile(w, r, "web/index.html")
	}
}
