package entity

import (
	"time"

	"github.com/google/uuid"
)

type InformationAndReferral struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	InfoType    string    `gorm:"size:100;not null" json:"infoType"`
	Description string    `gorm:"size:255" json:"description"`
	Date        time.Time `gorm:"type:date" json:"date"`

	PersonId uuid.UUID `gorm:"not null" json:"personId"`
}
