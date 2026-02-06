package entity

import "github.com/google/uuid"

type AddressPhone struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	PhoneNumber string    `gorm:"size:25;not null"`
	Type        string    `gorm:"size:50;not null"`

	PlaceId  uuid.UUID `gorm:"not null"`           // TODO: Embed
	Place    Place     `gorm:"foreignKey:PlaceId"` // TODO: Remove duplicate?
	PersonId uuid.UUID `gorm:"not null"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:AddressPhoneId"`
}
