package models

type ItemImage struct {
	ImageID   int    `gorm:"primaryKey;column:image_id;autoIncrement"`
	ImageData []byte `gorm:"not null;column:image_data"`
	IsPrimary bool   `gorm:"default:false;column:is_primary"`
	ItemID    int    `gorm:"not null;column:item_id"`
}

func (ItemImage) TableName() string {
	return "Item_image"
}
