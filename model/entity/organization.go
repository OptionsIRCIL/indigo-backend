package entity

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id   uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name string    `gorm:"size:150;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
