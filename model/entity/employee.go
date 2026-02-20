package entity

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	Id        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	FirstName string    `gorm:"size:255;not null"`
	LastName  string    `gorm:"size:255;not null"`
	Username  string    `gorm:"size:255;unique;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
