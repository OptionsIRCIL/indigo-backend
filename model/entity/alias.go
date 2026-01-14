package entity

type Alias struct {
	Id        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	AliasName string `gorm:"column:aliasName;size:100;not null" json:"aliasName"`

	PersonId uint `gorm:"column:personId;not null" json:"personId"`
}
