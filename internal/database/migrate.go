package database

import (
	"gorm.io/gorm"
	"log"
	"os"
)

func RunMigrations(db *gorm.DB) error {
	files := []string{
		"internal/database/migrations/001_create_users_table.sql",
		"internal/database/migrations/002_create_item_table.sql",
		"internal/database/migrations/003_create_itemImage_table.sql",
		"internal/database/migrations/004_create_purchase_table.sql",
		"internal/database/migrations/005_create_message_table.sql",
		"internal/database/migrations/006_create_chat_table.sql",
		"internal/database/migrations/007_create_dish_table.sql",
		"internal/database/migrations/008_alter_message_table.sql",
		"internal/database/migrations/009_create_forumBlockList_table.sql",
		"internal/database/migrations/010_alter_message_forumBlockList.sql",
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
