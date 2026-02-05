package entity

import "github.com/google/uuid"

type ConsumerService struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	ServiceName string    `gorm:"size:100;not null" json:"serviceName"`
	Status      string    `gorm:"size:50" json:"status"`

	PersonId uuid.UUID `gorm:"not null" json:"personId"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:ConsumerServiceId" json:"-"`
}
