package models

import "time"

type Message struct {
	MessageID   int64     `gorm:"column:message_id;primaryKey;autoIncrement:true"` // Соответствует BIGSERIAL
	ChatID      int       `gorm:"column:chat_id;not null"`                         // Связь с Chat
	SenderID    int       `gorm:"column:sender_id;not null"`                       // Связь с User
	Value       string    `gorm:"column:value;type:varchar(256);not null"`         // Текст сообщения
	MessageTime time.Time `gorm:"column:message_time;not null"`                    // Время отправки
}

func (Message) TableName() string {
	return "Message"
}
