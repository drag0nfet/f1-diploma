package models

type Hall struct {
	HallID      int      `gorm:"column:hall_id;primaryKey;autoIncrement" json:"hall_id"`
	TableGrid   [][]int  `gorm:"column:table_grid;type:int[][];not null" json:"table_grid"`
	Name        string   `gorm:"column:name;type:varchar(20);not null" json:"name"`
	Description string   `gorm:"column:description;type:varchar(256);not null" json:"description"`
	Album       [][]byte `gorm:"column:album;type:bytea[]" json:"album"`
	MaxTables   int      `gorm:"column:max_tables;not null" json:"max_tables"`
	NowTables   int      `gorm:"column:now_tables;not null;default:0" json:"now_tables"`
}

func (Hall) TableName() string {
	return `public."Hall"`
}
