package models

type Dish struct {
	DishID      int    `gorm:"primaryKey;column:dish_id;autoIncrement"` // serial -> autoIncrement
	Name        string `gorm:"not null;column:name;type:varchar"`       // varchar not null
	Cost        int    `gorm:"not null;column:cost"`                    // integer not null
	Description string `gorm:"column:description;type:text"`            // text (nullable by default)
	Image       []byte `gorm:"column:image;type:bytea"`                 // bytea (для хранения бинарных данных, например, изображений)
}

func (Dish) TableName() string {
	return "Dish"
}
