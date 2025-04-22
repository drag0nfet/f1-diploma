package models

import "time"

type News struct {
	NewsID      int       `gorm:"column:news_id;primaryKey;autoIncrement:true" json:"news_id"`
	CreatorID   int       `gorm:"column:creator_id;not null" json:"creator_id"`
	Title       string    `gorm:"column:title;type:varchar(64);not null" json:"title"`
	Description *string   `gorm:"column:description;type:varchar(256)" json:"description"`
	Comment     string    `gorm:"column:comment;type:text;not null" json:"comment"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Image       []byte    `gorm:"column:image;type:bytea" json:"image"`
	Status      string    `gorm:"column:status;type:varchar(10);default:DRAFT" json:"status"`
}

func (News) TableName() string {
	return "News"
}
