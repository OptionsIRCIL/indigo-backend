package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	Id                          uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Name                        string    `gorm:"size:150;not null"`
	Location                    string    `gorm:"size:150;not null"`
	LastContact                 time.Time `gorm:"type:datetime;not null"`
	Email                       string    `gorm:"size:255;not null"`
	Phone                       string    `gorm:"size:11;not null"`
	Address1                    string    `gorm:"size:150;not null"`
	Address2                    string    `gorm:"size:150;not null"`
	City                        string    `gorm:"size:150;not null"`
	State                       string    `gorm:"size:150;not null"`
	Zip                         string    `gorm:"size:150;not null"`
	County                      string    `gorm:"size:150;not null"`
	ContactFirstName            string    `gorm:"size:150;not null"`
	ContactLastName             string    `gorm:"size:150;not null"`
	ContactPosition             string    `gorm:"size:150;not null"`
	ContactDept                 string    `gorm:"size:150;not null"`
	IsCommunityService          bool      `gorm:"type:bool;not null"`
	OfferedServiceDescription   string    `gorm:"size:255;not null"`
	ServiceEligibleRequirements string    `gorm:"size:255;not null"`
	IntakeProcedures            string    `gorm:"size:255;not null"`
	FeesCharged                 string    `gorm:"size:255;not null"`
	HoursOfOperations           string    `gorm:"size:255;not null"`
	OtherInformation            string    `gorm:"size:255;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	o.Id, err = uuid.NewRandom()
	return err
}
