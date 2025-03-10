package models

import (
	"time"
)

type Purchase struct {
	PurchaseID       int        `gorm:"primaryKey;column:purchase_id;default:nextval('purchase_purchase_id_seq'::regclass)"`
	ItemID           int        `gorm:"not null;column:item_id"`
	BuyerID          int        `gorm:"not null;column:buyer_id"`
	PurchaseTime     *time.Time `gorm:"column:purchase_time"`
	PurchaseStatus   string     `gorm:"not null;column:purchase_status"`
	StatusTime       time.Time  `gorm:"not null;column:status_time"`
	PurchaseQuantity *int       `gorm:"column:purchase_quantity"`
}

func (Purchase) TableName() string {
	return "Purchase"
}
