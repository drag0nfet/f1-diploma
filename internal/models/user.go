package models

import "time"

type User struct {
	UserID            int       `gorm:"primaryKey;column:user_id;autoIncrement"`
	Login             string    `gorm:"unique;not null;column:login;type:text"`
	Email             string    `gorm:"unique;not null;column:email;default:''"`
	Password          string    `gorm:"not null;column:password"`
	Rights            int       `gorm:"not null;default:0;column:rights"`
	IsConfirmed       bool      `gorm:"column:is_confirmed;default:false"`
	ConfirmationToken string    `gorm:"column:confirmation_token"`
	LastSent          time.Time `gorm:"column:last_sent;default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "User"
}
