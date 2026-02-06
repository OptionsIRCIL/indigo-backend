package entity

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	GoalName    string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:255"`
	Status      string    `gorm:"size:50"`

	PersonId uuid.UUID `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
