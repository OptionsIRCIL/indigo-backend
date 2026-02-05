package entity

import "github.com/google/uuid"

type Alias struct {
	Id        uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	AliasName string    `gorm:"size:100;not null" json:"aliasName"`

	PersonId uuid.UUID `gorm:"not null" json:"personId"`
}
