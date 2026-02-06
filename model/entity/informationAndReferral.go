package entity

import (
	"time"

	"github.com/google/uuid"
)

type InformationAndReferral struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Date        time.Time `gorm:"type:date"`
	Hours       int64     `gorm:"type:bigint"`
	TravelHours int64     `gorm:"type:bigint"`
	Department  string    `gorm:"type:varchar(255)"`

	//Referral and Requests
	CallerType     string `gorm:"type:varchar(255)"`
	Disability     string `gorm:"type:varchar(255)"`
	Referrer       string `gorm:"type:varchar(255)"`
	ServiceRequest string `gorm:"type:varchar(255)"`
	Outcome        string `gorm:"type:varchar(255)"`

	//Logging
	FormDate    time.Time `gorm:"type:date"`
	ServiceType string    `gorm:"type:varchar(255)"`
	Grant       string    `gorm:"type:varchar(255)"`
	Units       string    `gorm:"type:varchar(255)"`

	EmployeeId     uuid.UUID `gorm:"type:char(36)"`
	OrganizationId uuid.UUID `gorm:"type:char(36)"`
	PersonId       uuid.UUID `gorm:"not null"`
}
