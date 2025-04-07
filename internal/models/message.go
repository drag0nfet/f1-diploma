package models

import "time"
import "gopkg.in/guregu/null.v4"

type Message struct {
	MessageID   int64     `gorm:"column:message_id;primaryKey;autoIncrement:true" json:"message_id"`
	ChatID      int       `gorm:"column:chat_id;not null" json:"chat_id"`
	SenderID    int       `gorm:"column:sender_id;not null" json:"sender_id"`
	Value       string    `gorm:"column:value;type:varchar(256);not null" json:"value"`
	MessageTime time.Time `gorm:"column:message_time;not null" json:"message_time"`
	ReplyID     null.Int  `gorm:"column:reply_id" json:"reply_id"`
}

func (Message) TableName() string {
	return "Message"
}
