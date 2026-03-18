package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	Id        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	FirstName string    `gorm:"size:255;not null" groups:"get"`
	LastName  string    `gorm:"size:255;not null" groups:"get"`
	Username  string    `gorm:"size:255;uniqueIndex;not null" groups:"get"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (e *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	e.Id, err = uuid.NewRandom()
	return err
}
