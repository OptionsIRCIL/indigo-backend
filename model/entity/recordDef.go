package entity

import (
	"time"

	"github.com/google/uuid"
)

type RecordDef struct {
	Id    uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Type  string    `gorm:"size:50;not null"`
	Value string    `gorm:"size:255;not null"`

	PersonId uuid.UUID `gorm:"column:personId;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
