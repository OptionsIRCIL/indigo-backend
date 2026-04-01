package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Grant struct {
	Id     uuid.UUID `gorm:"primaryKey; type: CHAR(36)" groups:"get"`
	Name   string    `gorm:"type: VARCHAR(255); not null; index" groups:"get,post"`
	Active bool
}

func (g *Grant) BeforeCreate(tx *gorm.DB) (err error) {
	g.Id, err = uuid.NewRandom()
	return err
}
