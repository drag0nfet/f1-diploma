package models

import "time"

type Chat struct {
	ChatID    int       `gorm:"column:chat_id;primaryKey;autoIncrement:true"`
	ChatType  string    `gorm:"column:chat_type;type:varchar(20);not null"`
	ItemID    *int      `gorm:"column:item_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	Title     *string   `gorm:"column:title;type:varchar(100)"`
}

func (Chat) TableName() string {
	return "Chat"
}
