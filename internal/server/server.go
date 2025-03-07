package server

import (
	"diploma/internal/handlers"
	"net/http"
)

func Run() {
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	http.ListenAndServe(":5050", nil)
}
