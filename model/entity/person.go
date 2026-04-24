package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/model"
)

type Person struct {
	Id            uuid.UUID  `gorm:"primaryKey; type:char(36)" groups:"get"`
	FirstName     string     `gorm:"size:100; not null" groups:"get,post" validate:"notblank"`
	LastName      string     `gorm:"size:100; not null" groups:"get,post" validate:"notblank"`
	Salutation    *string    `gorm:"size:20" groups:"get,post" default:"null"`
	Gender        *string    `gorm:"size:20" groups:"get,post" default:"null"`
	Ethnicity     string     `gorm:"size:100; not null" groups:"get,post" validate:"notblank"`
	OptNewsletter bool       `gorm:"not null; default:true" groups:"get,post" default:"true"`
	Birthday      model.Date `groups:"get,post"`
	Email         string     `gorm:"size:150; unique; not null" groups:"get,post" validate:"notblank,email"`
	Phone         string     `gorm:"size:25; not null" groups:"get,post" validate:"notblank,phone"`
	Active        bool       `gorm:"default:true; not null" groups:"get,post" default:"true"`
	Deceased      bool       `gorm:"default:false; not null" groups:"get,post" default:"false"`
	Disabilities  string     `gorm:"type: TEXT(5000)" groups:"get,post"`

	AddressPhones           []AddressPhone           `gorm:"foreignKey:PersonId"`
	Aliases                 []PersonAlias            `gorm:"foreignKey:PersonId"`
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
