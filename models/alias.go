package models

type Alias struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	AliasName string `gorm:"column:aliasName;size:100;not null" json:"aliasName"`

	PersonID uint `gorm:"column:personId;not null" json:"personId"`
}
