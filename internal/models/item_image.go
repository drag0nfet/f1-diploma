package models

type ItemImage struct {
	ImageID   int    `gorm:"primaryKey;column:image_id;autoIncrement"` // Primary key
	ImageData []byte `gorm:"not null;column:image_data"`               // Binary data for the image
	IsPrimary bool   `gorm:"default:false;column:is_primary"`          // Is this the primary image
	ItemID    int    `gorm:"not null;column:item_id"`                  // Foreign key to the item
}

func (ItemImage) TableName() string {
	return "Item_image"
}
