package models

import "time"

type BookingSpot struct {
	BookingID int       `gorm:"column:booking_id;primaryKey;autoIncrement" json:"booking_id"`
	SpotID    int       `gorm:"column:spot_id;not null" json:"spot_id"`
	UserID    *int      `gorm:"column:user_id" json:"user_id"`
	EventID   *int      `gorm:"column:event_id" json:"event_id"`
	Status    string    `gorm:"column:status;type:varchar(10);not null;default:INACTIVE" json:"status"`
	StartTime time.Time `gorm:"column:start_time;type:timestamptz;not null" json:"start_time"`
}

func (BookingSpot) TableName() string {
	return `public."BookingSpot"`
}
