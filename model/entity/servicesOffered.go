package entity

import (
	"time"

	"github.com/google/uuid"
)

type ServicesOffered struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	ServiceName string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:255"`

	ConsumerServiceId uuid.UUID `gorm:"not null"`
	AddressPhoneId    uuid.UUID `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
