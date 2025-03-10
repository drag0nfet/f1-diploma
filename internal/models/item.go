package models

type Item struct {
	ItemID      int     `gorm:"primaryKey;column:item_id;autoIncrement"`
	ItemName    string  `gorm:"not null;column:name"`
	Description string  `gorm:"type:text;column:description"`
	Price       float64 `gorm:"not null;column:price"`
	UserID      int     `gorm:"not null;column:user_id"`
}

func (Item) TableName() string {
	return "Item"
}
