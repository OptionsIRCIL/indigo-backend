package entity

import (
	"time"

	"github.com/google/uuid"
)

type CommunityServiceEvent struct {
	Id                      uuid.UUID `gorm:"primaryKey;type:char(36)"`
	InitialDate             time.Time `gorm:"type:datetime;not null"`
	Category                string    `gorm:"size:100;not null"` // TODO: Enum?
	FutureReference         string    `gorm:"size:255"`          // TODO: Remove or rename
	ServiceDescription      string    `gorm:"size:255;not null"` // TODO: TEXT type?
	Outcome                 string    `gorm:"size:255"`
	Publications            int       `gorm:"default:0"`
	PersonsWithDisabilities int       `gorm:"default:0"`
	GeneralPublic           int       `gorm:"default:0"`
}
