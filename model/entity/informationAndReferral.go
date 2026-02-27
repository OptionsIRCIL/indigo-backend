package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/model"
)

type InformationAndReferral struct {
	Id         uuid.UUID  `gorm:"primaryKey;type:char(36)" groups:"get"`
	Date       model.Date `groups:"get,post"`
	Hours      int64      `gorm:"type:bigint" groups:"get,post"`
	Department string     `gorm:"type:varchar(255)" groups:"get,post"`

	//Referral and Requests
	Referrer       string `gorm:"type:varchar(255)" groups:"get,post"`
	ServiceRequest string `gorm:"type:varchar(255)" groups:"get,post"`
	Outcome        string `gorm:"type:varchar(255)" groups:"get,post"`

	//Logging (Re-Use?)
	FormDate    model.Date `groups:"get,post"`
	ServiceType string     `gorm:"type:varchar(255)" groups:"get,post"`
	Grant       string     `gorm:"type:varchar(255)" groups:"get,post"`

	EmployeeId     uuid.UUID `gorm:"type:char(36)" groups:"get,post"`
	OrganizationId uuid.UUID `gorm:"type:char(36)" groups:"get,post"`
	PersonId       uuid.UUID `gorm:"type:char(36); not null" groups:"get,post"`

	CreatedAt time.Time `groups:"get"`
	UpdatedAt time.Time `groups:"get"`
	DeletedAt time.Time
}

func (i *InformationAndReferral) BeforeCreate(tx *gorm.DB) (err error) {
	i.Id, err = uuid.NewRandom()
	return err
}
