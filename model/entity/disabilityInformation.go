package entity

import (
	"time"

	"github.com/google/uuid"
)

type DisabilityInfo struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Disability  string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:255"`
	Severity    string    `gorm:"size:50"`

	PersonId uuid.UUID `gorm:"unique;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
