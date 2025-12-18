package models

import "time"

type InformationAndReferral struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InfoType    string    `gorm:"column:infoType;size:100;not null" json:"infoType"`
	Description string    `gorm:"column:description;size:255" json:"description"`
	Date        time.Time `gorm:"column:date;type:date" json:"date"`

	PersonID uint `gorm:"column:personId;not null" json:"personId"`
}
