package main

import (
	"diploma/internal/database"
	"diploma/internal/server"
	"log"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatal("Ошибка при подключении к БД:", err)
	}

	server.Run()
}
