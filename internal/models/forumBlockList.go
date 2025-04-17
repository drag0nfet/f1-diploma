package models

import "time"

type ForumBlockList struct {
	UserID      int       `gorm:"column:user_id;not null" json:"user_id"`
	MessageID   int       `gorm:"column:message_id;not null" json:"message_id"`
	ModeratorID int       `gorm:"column:moderator_id;not null" json:"moderator_id"`
	IsValid     bool      `gorm:"column:is_valid;default:true" json:"is_valid"`
	TimeGot     time.Time `gorm:"column:time_got;not null" json:"time_got"`
	Status      string    `gorm:"column:status;type:varchar(20);not null;default:READY" json:"status"`
}

func (ForumBlockList) TableName() string {
	return "ForumBlockList"
}
