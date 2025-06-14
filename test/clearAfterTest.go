package test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	dbpass := os.Getenv("DB_PASS")
	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=f1-diploma sslmode=disable", dbpass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// 1 Удаляем тестовые сообщения с value = 'ТЕСТОВОЕ_SMS'
	if _, err := db.Exec(`DELETE FROM "Message" WHERE value = 'ТЕСТОВОЕ_SMS'`); err != nil {
		log.Println("Ошибка удаления сообщений:", err.Error())
	} else {
		fmt.Println("Тестовые сообщения удалены.")
	}

	// 2 бронирования

	query := `
		DELETE FROM "BookingSpot"
		WHERE user_id>999
	`

	if _, err := db.Exec(query); err != nil {
		log.Println("Ошибка удаления бронирований:", err)
	} else {
		fmt.Println("Тестовые бронирования удалены.")
	}

	// 3 Удаляем тестовых пользователей с login, начинающимся с K6TEST_
	if _, err := db.Exec(`DELETE FROM "User" WHERE login LIKE 'K6TEST_%'`); err != nil {
		log.Println("Ошибка удаления пользователей:", err)
	} else {
		fmt.Println("Тестовые пользователи удалены.")
	}

}
