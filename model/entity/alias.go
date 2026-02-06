package entity

import "github.com/google/uuid"

type Alias struct {
	AliasName string    `gorm:"primaryKey;size:100;not null"`
	PersonId  uuid.UUID `gorm:"primaryKey;not null"`
}
