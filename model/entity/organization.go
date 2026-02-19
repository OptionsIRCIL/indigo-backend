package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	Id   uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name string    `gorm:"size:150;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	o.Id, err = uuid.NewRandom()
	return err
}
