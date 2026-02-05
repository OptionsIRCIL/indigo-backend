package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	Id         uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	FirstName  string    `gorm:"size:100;not null" json:"firstName"`
	LastName   string    `gorm:"size:100;not null" json:"lastName"`
	Salutation string    `gorm:"size:20" json:"salutation"`
	Gender     string    `gorm:"size:20" json:"gender"`
	Birthday   time.Time `gorm:"type:date" json:"birthday"`
	Email      string    `gorm:"size:150;unique;not null" json:"email"`
	Phone      string    `gorm:"size:25" json:"phone"`
	Active     bool      `gorm:"default:true" json:"active"`
	Deceased   bool      `gorm:"default:false" json:"deceased"`

	AddressPhones           []AddressPhone           `gorm:"foreignKey:PersonId" json:"-"`
	Aliases                 []Alias                  `gorm:"foreignKey:PersonId" json:"-"`
	RecordDefs              []RecordDef              `gorm:"foreignKey:PersonId" json:"-"`
	DisabilityInfo          *DisabilityInfo          `gorm:"foreignKey:PersonId" json:"-"`
	InformationAndReferrals []InformationAndReferral `gorm:"foreignKey:PersonId" json:"-"`
	ConsumerServices        []ConsumerService        `gorm:"foreignKey:PersonId" json:"-"`
	Goals                   []Goal                   `gorm:"foreignKey:PersonId" json:"-"`
}

func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id, err = uuid.NewRandom()
	return err
}
