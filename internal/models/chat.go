package models

import "time"

type Chat struct {
	ChatID    int       `gorm:"column:chat_id;primaryKey;autoIncrement:true"` // Соответствует SERIAL
	ChatType  string    `gorm:"column:chat_type;type:varchar(20);not null"`   // Тип чата
	ItemID    *int      `gorm:"column:item_id"`                               // Связь с Item (может быть NULL)
	CreatedAt time.Time `gorm:"column:created_at;not null"`                   // Время создания
	Title     *string   `gorm:"column:title;type:varchar(100)"`               // Название чата (может быть NULL)
}

func (Chat) TableName() string {
	return "Chat"
}
