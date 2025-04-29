package database

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	migrationsDir := "internal/database/migrations"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	var files []fs.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		path := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err = db.Exec(string(content)).Error; err != nil {
			return err
		}
	}

	log.Println("Миграция БД успешна")
	return nil
}
