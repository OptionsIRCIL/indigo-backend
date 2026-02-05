package entity

import "github.com/google/uuid"

type ServicesOffered struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	ServiceName string    `gorm:"size:100;not null" json:"serviceName"`
	Description string    `gorm:"size:255" json:"description"`

	ConsumerServiceId uuid.UUID `gorm:"not null" json:"consumerServiceId"`
	AddressPhoneId    uuid.UUID `gorm:"not null" json:"addressPhoneId"`
}
