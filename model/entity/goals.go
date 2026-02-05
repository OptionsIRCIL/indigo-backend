package entity

import "github.com/google/uuid"

type Goal struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	GoalName    string    `gorm:"size:100;not null" json:"goalName"`
	Description string    `gorm:"size:255" json:"description"`
	Status      string    `gorm:"size:50" json:"status"`

	PersonId uuid.UUID `gorm:"not null" json:"personId"`
}
