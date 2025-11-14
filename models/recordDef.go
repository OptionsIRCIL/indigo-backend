package models

type RecordDef struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Type  string `gorm:"column:type;size:50;not null" json:"type"`
	Value string `gorm:"column:value;size:255;not null" json:"value"`

	PersonID uint `gorm:"column:personId;not null" json:"personId"`
}
