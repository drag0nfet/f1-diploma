package models

type User struct {
	UserID   int    `gorm:"primaryKey;column:user_id;autoIncrement"`
	Login    string `gorm:"unique;not null;column:login;type:text"`
	Password string `gorm:"not null;column:password"`
	Rights   int    `gorm:"not null;default:0;column:rights"`
}

func (User) TableName() string {
	return "User"
}
