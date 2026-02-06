package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	Id         uuid.UUID `gorm:"primaryKey;type:char(36)" groups:"get"`
	FirstName  string    `gorm:"size:100;not null" groups:"get,post"`
	LastName   string    `gorm:"size:100;not null" groups:"get,post"`
	Salutation string    `gorm:"size:20" groups:"get,post"`
	Gender     string    `gorm:"size:20" groups:"get,post"`
	Birthday   time.Time `gorm:"type:date" groups:"get,post"`
	Email      string    `gorm:"size:150;unique;not null" groups:"get,post"`
	Phone      string    `gorm:"size:25" groups:"get,post"`
	Active     bool      `gorm:"default:true" groups:"get,post"`
	Deceased   bool      `gorm:"default:false" groups:"get,post"`

	AddressPhones           []AddressPhone           `gorm:"foreignKey:PersonId"`
	Aliases                 []Alias                  `gorm:"foreignKey:PersonId"`
	RecordDefs              []RecordDef              `gorm:"foreignKey:PersonId"`
	DisabilityInfo          *DisabilityInfo          `gorm:"foreignKey:PersonId"`
	InformationAndReferrals []InformationAndReferral `gorm:"foreignKey:PersonId"`
	ConsumerServices        []ConsumerService        `gorm:"foreignKey:PersonId"`
	Goals                   []Goal                   `gorm:"foreignKey:PersonId"`

	CreatedAt time.Time `groups:"get"`
	UpdatedAt time.Time `groups:"get"`
	DeletedAt time.Time
}

func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id, err = uuid.NewRandom()
	return err
}
