package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Identifier struct {
	Id string `gorm:"primaryKey;size:36" json:"id"`
}

// BeforeCreate is GORM Hook that runs automatically before item should be added to the database
func (i *Identifier) BeforeCreate(tx *gorm.DB) (err error) {
	// Only generate a UUID if one wasn't manually set
	if i.Id == "" {
		i.Id = uuid.New().String()
	}
	return
}
