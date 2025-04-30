package models

type Spot struct {
	SpotID   int `gorm:"column:spot_id;primaryKey;autoIncrement" json:"spot_id"`
	TableID  int `gorm:"column:table_id;not null" json:"table_id"`
	SpotName int `gorm:"column:spot_name;not null" json:"spot_name"`
}

func (Spot) TableName() string {
	return `public."Spot"`
}
