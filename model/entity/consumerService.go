package entity

import "github.com/google/uuid"

type ConsumerService struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	ServiceName string    `gorm:"size:100;not null"`
	Status      string    `gorm:"size:50"`

	PersonId uuid.UUID `gorm:"not null"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:ConsumerServiceId"`
}
