package entity

import "time"

type InformationAndReferral struct {
	Id          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InfoType    string    `gorm:"column:infoType;size:100;not null" json:"infoType"`
	Description string    `gorm:"column:description;size:255" json:"description"`
	Date        time.Time `gorm:"column:date;type:date" json:"date"`

	PersonId uint `gorm:"column:personId;not null" json:"personId"`
}
