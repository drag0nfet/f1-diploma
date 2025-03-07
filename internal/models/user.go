package models

type User struct {
	UserID   int    `gorm:"primaryKey;autoIncrement;column:user_id"`
	Login    string `gorm:"uniqueIndex;not null;column:login"`
	Password string `gorm:"not null;column:password"`
	Rights   int    `gorm:"not null;default:0;column:rights"`
}

func (User) TableName() string {
	return "User"
}
