package models

type RecordDef struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Type  string `gorm:"column:type;size:50;not null"`
	Value string `gorm:"column:value;size:255;not null"`

	PersonID uint `gorm:"column:personId;not null"`
}
