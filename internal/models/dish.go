package models

type Dish struct {
	DishID      int    `gorm:"primaryKey;column:dish_id;autoIncrement"`
	Name        string `gorm:"not null;column:name;type:varchar"`
	Cost        int    `gorm:"not null;column:cost"`
	Description string `gorm:"column:description;type:text"`
	Image       []byte `gorm:"column:image;type:bytea"`
}

func (Dish) TableName() string {
	return "Dish"
}
