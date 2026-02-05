package entity

import "github.com/google/uuid"

type DisabilityInfo struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	Disability  string    `gorm:"size:100;not null" json:"disability"`
	Description string    `gorm:"size:255" json:"description"`
	Severity    string    `gorm:"size:50" json:"severity"`

	PersonId uuid.UUID `gorm:"unique;not null" json:"personId"`
}
