package entity

import (
	"time"

	"github.com/google/uuid"
)

type InformationAndReferral struct {
	Id          uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	Date        time.Time `gorm:"type:date" json:"date"`
	Hours       int64     `gorm:"type:bigint" json:"hours"`
	TravelHours int64     `gorm:"type:bigint" json:"travelHours"`
	Department  string    `gorm:"type:varchar(255)" json:"department"`

	//Personal Intake Information
	Name      string `gorm:"type:varchar(255)" json:"name"`
	Address   string `gorm:"type:varchar(255)" json:"address"`
	City      string `gorm:"type:varchar(255)" json:"city"`
	State     string `gorm:"type:varchar(255)" json:"state"`
	Zip       int64  `gorm:"type:bigint" json:"zip"`
	Phone     string `gorm:"type:varchar(255)" json:"phone"`
	County    string `gorm:"type:varchar(255)" json:"county"`
	Email     string `gorm:"type:varchar(255)" json:"email"`
	Gender    string `gorm:"type:varchar(255)" json:"gender"`
	Ethnicity string `gorm:"type:varchar(255)" json:"ethnicity"`

	//Referral and Requests
	CallerType     string `gorm:"type:varchar(255)" json:"callerType"`
	Disability     string `gorm:"type:varchar(255)" json:"disability"`
	Referrer       string `gorm:"type:varchar(255)" json:"referrer"`
	ServiceRequest string `gorm:"type:varchar(255)" json:"serviceRequest"`
	Outcome        string `gorm:"type:varchar(255)" json:"outcome"`

	//Logging
	FormDate    time.Time `gorm:"type:date" json:"formDate"`
	ServiceType string    `gorm:"type:varchar(255)" json:"serviceType"`
	Grant       string    `gorm:"type:varchar(255)" json:"grant"`
	Units       string    `gorm:"type:varchar(255)" json:"units"`

	EmployeeId     uuid.UUID `gorm:"type:char(36)" json:"employeeId"`
	OrganizationId uuid.UUID `gorm:"type:char(36)" json:"organizationId"`
	PersonId       uuid.UUID `gorm:"not null" json:"personId"`
}
