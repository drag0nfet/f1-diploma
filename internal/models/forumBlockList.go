package models

type ForumBlockList struct {
	UserID      int  `gorm:"column:user_id;not null" json:"user_id"`
	MessageID   int  `gorm:"column:message_id;not null" json:"message_id"`
	ModeratorID int  `gorm:"column:moderator_id;not null" json:"moderator_id"`
	IsValid     bool `gorm:"column:is_valid;default:true" json:"is_valid"`
}

func (ForumBlockList) TableName() string {
	return "ForumBlockList"
}
