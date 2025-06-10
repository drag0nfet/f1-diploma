package main

import (
	"diploma/internal/database"
	"diploma/internal/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	err = database.InitDB()
	if err != nil {
		log.Fatal("Ошибка при подключении к БД:", err)
	}

	server.Run()
}
