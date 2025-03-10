package server

import (
	"diploma/internal/handlers"
	"log"
	"net/http"
)

func Run() {
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	err := http.ListenAndServe("192.168.30.11:5051", nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	} else {
		log.Println("ListenAndServe success")
	}
}
