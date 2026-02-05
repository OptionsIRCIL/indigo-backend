package entity

import (
	"time"

	"github.com/google/uuid"
)

type CommunityServiceEvent struct {
	Id                      uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	InitialDate             time.Time `gorm:"type:datetime;not null" json:"initialDate"`
	Category                string    `gorm:"size:100;not null" json:"category"`           // TODO: Enum?
	FutureReference         string    `gorm:"size:255" json:"futureReference"`             // TODO: Remove or rename
	ServiceDescription      string    `gorm:"size:255;not null" json:"serviceDescription"` // TODO: TEXT type?
	Outcome                 string    `gorm:"size:255" json:"outcome"`
	Publications            int       `gorm:"default:0" json:"publications"`
	PersonsWithDisabilities int       `gorm:"default:0" json:"personsWithDisabilities"`
	GeneralPublic           int       `gorm:"default:0" json:"generalPublic"`
}
