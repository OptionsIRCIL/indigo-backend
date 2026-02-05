package entity

import "github.com/google/uuid"

type RecordDef struct {
	Id    uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	Type  string    `gorm:"size:50;not null" json:"type"`
	Value string    `gorm:"size:255;not null" json:"value"`

	PersonId uuid.UUID `gorm:"column:personId;not null" json:"personId"`
}
