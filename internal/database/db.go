package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "host=localhost user=postgres password=njseL]:u%!ZUc;2Y dbname=f1-diploma port=5432 sslmode=disable"
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
