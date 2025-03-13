package server

import (
	"diploma/internal/handlers"
	"diploma/internal/handlers/account"
	"diploma/internal/handlers/index"
	"diploma/internal/services"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web"))

	// Обработчик для статических файлов (JS, CSS)
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
		index.Register(w, r)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		index.Login(w, r)
	})
	mux.HandleFunc("/logout", account.Logout)
	mux.HandleFunc("/check-auth", handlers.CheckAuth)

	// Страничные маршруты
	mux.HandleFunc("/", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	}))
	mux.HandleFunc("/account", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/account.html")
	}))

	handler := services.EnableCORS(mux)

	err := http.ListenAndServe(SetIP()+":5051", handler)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	} else {
		log.Println("ListenAndServe success")
	}
}
