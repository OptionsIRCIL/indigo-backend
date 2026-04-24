package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityEducationAndOutreach struct {
	Id uuid.UUID `gorm:"primaryKey; type: CHAR(36)" groups:"get"`

	Category []string

	CreatedAt time.Time `groups:"get"`
	UpdatedAt time.Time `groups:"get"`
	DeletedAt time.Time
}

func (i *CommunityEducationAndOutreach) BeforeCreate(tx *gorm.DB) (err error) {
	i.Id, err = uuid.NewRandom()
	return err
}
