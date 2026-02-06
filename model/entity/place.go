package entity

import (
	"github.com/google/uuid"
)

type Place struct {
	Id           uuid.UUID `gorm:"primaryKey;type:char(36)"`
	AddressLine1 string    `gorm:"size:255;not null"`
	AddressLine2 string    `gorm:"size:255"`
	City         string    `gorm:"size:100;not null"`
	State        string    `gorm:"size:100;not null"`
	ZipCode      string    `gorm:"size:20"`
	Country      string    `gorm:"size:100;not null"`
}
