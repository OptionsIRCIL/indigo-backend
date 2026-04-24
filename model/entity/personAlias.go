package entity

import "github.com/google/uuid"

type PersonAlias struct {
	FirstName string    `gorm:"primaryKey;size:100;not null"`
	LastName  string    `gorm:"primaryKey;size:100;not null"`
	PersonId  uuid.UUID `gorm:"primaryKey;not null"`
}
