package entity

import "github.com/google/uuid"

type AddressPhone struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	PhoneNumber string    `gorm:"size:25;not null" json:"phoneNumber"`
	Type        string    `gorm:"size:50;not null" json:"type"`

	PlaceId  uuid.UUID `gorm:"not null" json:"placeId"`     // TODO: Embed
	Place    Place     `gorm:"foreignKey:PlaceId" json:"-"` // TODO: Remove duplicate?
	PersonId uuid.UUID `gorm:"not null" json:"personId"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:AddressPhoneId" json:"-"`
}
