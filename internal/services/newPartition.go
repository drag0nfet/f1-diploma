package services

import (
	"diploma/internal/database"
	"fmt"
	"log"
)

func NewPartition(id int64) {
	partitionNum := id / 700000

	startID := 700000 * (partitionNum + 1)
	endID := 700000 * (partitionNum + 2)

	partitionName := fmt.Sprintf("message_%d_%d", startID, endID)

	var exists bool
	checkSQL := `SELECT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = $1)`
	if err := database.DB.Raw(checkSQL, partitionName).Scan(&exists).Error; err != nil {
		log.Printf("Ошибка при проверке существования партиции %s: %v", partitionName, err)
		return
	}

	if exists {
		log.Printf("Партиция %s уже существует", partitionName)
		return
	}

	tx := database.DB.Begin()
	createPartitionSQL := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s
            PARTITION OF "Message"
            (
                CONSTRAINT "Message_sender_id_fkey_%d"
                    FOREIGN KEY (sender_id) REFERENCES "User",
                CONSTRAINT "Message_reply_id_fkey_%d"
                    FOREIGN KEY (reply_id) REFERENCES "Message"
                    ON DELETE SET NULL
            )
            FOR VALUES FROM (%d) TO (%d);
        `, partitionName, partitionNum+1, partitionNum+1, startID, endID)

	if err := tx.Exec(createPartitionSQL).Error; err != nil {
		tx.Rollback()
		log.Printf("Ошибка при создании партиции %s: %v", partitionName, err)
		return
	}

	alterOwnerSQL := fmt.Sprintf(`ALTER TABLE %s OWNER TO postgres;`, partitionName)
	if err := tx.Exec(alterOwnerSQL).Error; err != nil {
		tx.Rollback()
		log.Printf("Ошибка при установке владельца для партиции %s: %v", partitionName, err)
		return
	}

	createIndexSQL := fmt.Sprintf(`
            CREATE INDEX IF NOT EXISTS idx_chat_id_time_%d
            ON %s (chat_id, message_time);
        `, partitionNum+1, partitionName)
	if err := tx.Exec(createIndexSQL).Error; err != nil {
		tx.Rollback()
		log.Printf("Ошибка при создании индекса для партиции %s: %v", partitionName, err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Ошибка при фиксации транзакции для партиции %s: %v", partitionName, err)
		return
	}

	log.Printf("Создана новая партиция %s для message_id от %d до %d", partitionName, startID, endID)
}
