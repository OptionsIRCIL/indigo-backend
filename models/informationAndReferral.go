package models

import "time"

type InformationAndReferral struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	InfoType    string    `gorm:"column:infoType;size:100;not null"`
	Description string    `gorm:"column:description;size:255"`
	Date        time.Time `gorm:"column:date;type:date"`

	PersonID uint `gorm:"column:personId;not null"`
}
