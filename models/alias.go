package models

type Alias struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	AliasName string `gorm:"column:aliasName;size:100;not null"`

	PersonID uint `gorm:"column:personId;not null"`
}
