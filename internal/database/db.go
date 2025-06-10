package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() error {

	dbpass := os.Getenv("DB_PASS")
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=f1-diploma port=5432 sslmode=disable", dbpass)
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Println("Соединение с БД успешно")

	err = RunMigrations(DB)
	if err != nil {
		return err
	}

	log.Println("Подключение к БД успешно")
	return nil
}
