package models

import "time"

// Hall представляет зал без поля для фотографий, так как они хранятся в HallPhotos
type Hall struct {
	HallID      int    `gorm:"column:hall_id;primaryKey;autoIncrement" json:"hall_id"`
	Name        string `gorm:"column:name;type:varchar(20);not null" json:"name"`
	Description string `gorm:"column:description;type:varchar(256);not null" json:"description"`
}

func (Hall) TableName() string {
	return `public."Hall"`
}

// HallPhoto представляет фотографию зала
type HallPhoto struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	HallID    int       `gorm:"column:hall_id;not null" json:"hall_id"`
	Content   []byte    `gorm:"column:content;type:bytea;not null" json:"content"`
	MimeType  string    `gorm:"column:mime_type;type:varchar(50);not null" json:"mime_type"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp" json:"created_at"`
}

func (HallPhoto) TableName() string {
	return `public."HallPhotos"`
}
