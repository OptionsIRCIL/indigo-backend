package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InformationAndReferralEffort struct {
	Id uuid.UUID `gorm:"primaryKey; type: CHAR(36)" groups:"get"`

	// The employee tracking time
	EmployeeId uuid.UUID `gorm:"type: CHAR(36); not null" groups:"get,post"`
	Employee   Employee

	// The I&R time is being added to
	InformationAndReferralId uuid.UUID `gorm:"type: CHAR(36); not null" groups:"get,post"`
	InformationAndReferral   InformationAndReferral

	// The grant covering employee effort
	GrantId uuid.UUID `gorm:"type: CHAR(36); not null" groups:"get,post"`
	Grant   Grant

	// The number of minutes to add. For example, 2 hours of work would be reported as `{"minutes": 120}`.
	Minutes uint `gorm:"type: INTEGER UNSIGNED; not null" groups:"get,post"`

	CreatedAt time.Time `gorm:"not null" groups:"get"`
	UpdatedAt time.Time `gorm:"not null" groups:"get"`
	DeletedAt time.Time
}

func (i *InformationAndReferralEffort) BeforeCreate(tx *gorm.DB) (err error) {
	i.Id, err = uuid.NewRandom()
	return err
}
