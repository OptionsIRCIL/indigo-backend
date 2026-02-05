package entity

import (
	"github.com/google/uuid"
)

type Place struct {
	Id           uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	AddressLine1 string    `gorm:"size:255;not null" json:"addressLine1"`
	AddressLine2 string    `gorm:"size:255" json:"addressLine2"`
	City         string    `gorm:"size:100;not null" json:"city"`
	State        string    `gorm:"size:100;not null" json:"state"`
	ZipCode      string    `gorm:"size:20" json:"zipCode"`
	Country      string    `gorm:"size:100;not null" json:"country"`
}
