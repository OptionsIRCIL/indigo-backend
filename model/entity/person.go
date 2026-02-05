package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	Id         uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id" groups:"get"`
	FirstName  string    `gorm:"size:100;not null" json:"firstName" groups:"get,post"`
	LastName   string    `gorm:"size:100;not null" json:"lastName" groups:"get,post"`
	Salutation string    `gorm:"size:20" json:"salutation" groups:"get,post"`
	Gender     string    `gorm:"size:20" json:"gender" groups:"get,post"`
	Birthday   time.Time `gorm:"type:date" json:"birthday" groups:"get,post"`
	Email      string    `gorm:"size:150;unique;not null" json:"email" groups:"get,post"`
	Phone      string    `gorm:"size:25" json:"phone" groups:"get,post"`
	Active     bool      `gorm:"default:true" json:"active" groups:"get,post"`
	Deceased   bool      `gorm:"default:false" json:"deceased" groups:"get,post"`

	AddressPhones           []AddressPhone           `gorm:"foreignKey:PersonId" json:"-"`
	Aliases                 []Alias                  `gorm:"foreignKey:PersonId" json:"-"`
	RecordDefs              []RecordDef              `gorm:"foreignKey:PersonId" json:"-"`
	DisabilityInfo          *DisabilityInfo          `gorm:"foreignKey:PersonId" json:"-"`
	InformationAndReferrals []InformationAndReferral `gorm:"foreignKey:PersonId" json:"-"`
	ConsumerServices        []ConsumerService        `gorm:"foreignKey:PersonId" json:"-"`
	Goals                   []Goal                   `gorm:"foreignKey:PersonId" json:"-"`

	CreatedAt time.Time `groups:"get"`
	UpdatedAt time.Time `groups:"get"`
	DeletedAt time.Time
}

func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id, err = uuid.NewRandom()
	return err
}
