package models

type Table struct {
	TableID     int    `gorm:"column:table_id;primaryKey;autoIncrement" json:"table_id"`
	HallID      int    `gorm:"column:hall_id;not null" json:"hall_id"`
	TableNamee  int    `gorm:"column:table_name;not null" json:"table_name"`
	PriceStatus string `gorm:"column:price_status;type:varchar(20);not null" json:"price_status"`
	Seats       int    `gorm:"column:seats;not null" json:"seats"`
}

func (Table) TableName() string {
	return `public."Table"`
}
