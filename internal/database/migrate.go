package database

import (
	"gorm.io/gorm"
	"log"
	"os"
)

func RunMigrations(db *gorm.DB) error {
	files := []string{
		"internal/database/migrations/001_create_users_table.sql",
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if err = db.Exec(string(content)).Error; err != nil {
			return err
		}
	}
	log.Println("Миграция БД успешно выполнена.")
	return nil
}
