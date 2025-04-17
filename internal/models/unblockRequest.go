package models

import "time"

type UnblockRequest struct {
	RequestID int       `gorm:"column:request_id;primaryKey;autoIncrement:true" json:"request_id"`     // Соответствует SERIAL
	UserID    int       `gorm:"column:user_id;not null" json:"user_id"`                                // Пользователь, связанный с блокировкой
	MessageID int       `gorm:"column:message_id;not null" json:"message_id"`                          // Сообщение, связанное с блокировкой
	Status    string    `gorm:"column:status;type:varchar(20);not null;default:PENDING" json:"status"` // Статус запроса
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`                          // Время создания
	Comment   string    `gorm:"column:comment;type:text;not null" json:"comment"`                      // Комментарий
}

func (UnblockRequest) TableName() string {
	return "UnblockRequest"
}
