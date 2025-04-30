package models

import "time"

type Event struct {
	EventID       int       `gorm:"column:event_id;primaryKey;autoIncrement" json:"event_id"`
	Description   string    `gorm:"column:description;type:text;not null" json:"description"`
	TimeStart     time.Time `gorm:"column:time_start;type:timestamptz;not null" json:"time_start"`
	SportCategory string    `gorm:"column:sport_category;type:varchar(32);not null" json:"sport_category"`
	SportType     string    `gorm:"column:sport_type;type:varchar(32);not null" json:"sport_type"`
	PriceStatus   string    `gorm:"column:price_status;type:varchar(20);not null" json:"price_status"`
	Duration      int       `gorm:"column:duration;not null;default:90" json:"duration"`
}

func (Event) TableName() string {
	return `public."Event"`
}
