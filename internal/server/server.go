package server

import (
	"diploma/internal/handlers"
	"diploma/internal/services"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web"))

	// JS & CSS настройка
	mux.HandleFunc("/web/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/web/styles.css" {
			w.Header().Set("Content-Type", "text/css")
		}
		if r.URL.Path == "/web/js/main.js" {
			w.Header().Set("Content-Type", "application/javascript")
		}
		http.StripPrefix("/web/", fileServer).ServeHTTP(w, r)
	})

	// API маршруты
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(w, r)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(w, r)
	})

	// Системный маршрут для проверки авторизации
	mux.HandleFunc("/check-auth", handlers.CheckAuth)

	// Страничные маршруты
	mux.HandleFunc("/", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.Index(w, r)
	}))

	handler := services.EnableCORS(mux)

	err := http.ListenAndServe(SetIP()+":5051", handler)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	} else {
		log.Println("ListenAndServe success")
	}
}
