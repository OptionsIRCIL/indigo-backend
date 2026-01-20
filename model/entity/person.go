package entity

import "time"

type Person struct {
	Id         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName  string    `gorm:"column:firstName;size:100;not null" json:"firstName" validate:"required"`
	LastName   string    `gorm:"column:lastName;size:100;not null" json:"lastName" validate:"required"`
	Salutation string    `gorm:"column:salutation;size:20" json:"salutation" validate:"required"`
	Gender     string    `gorm:"column:gender;size:20" json:"gender" validate:"required"`
	Birthday   time.Time `gorm:"column:birthday;type:date" json:"birthday" validate:"required"`
	Email      string    `gorm:"column:email;size:150;unique;not null" json:"email" validate:"required"`
	Phone      string    `gorm:"column:phone;size:25" json:"phone" validate:"required"`
	Active     bool      `gorm:"column:active;default:true" json:"active" validate:"required"`
	Deceased   bool      `gorm:"column:deceased;default:false" json:"deceased" validate:"required"`

	AddressPhones           []AddressPhone           `gorm:"foreignKey:PersonId" json:"-"`
	Aliases                 []Alias                  `gorm:"foreignKey:PersonId" json:"-"`
	RecordDefs              []RecordDef              `gorm:"foreignKey:PersonId" json:"-"`
	DisabilityInfo          *DisabilityInfo          `gorm:"foreignKey:PersonId" json:"-"`
	InformationAndReferrals []InformationAndReferral `gorm:"foreignKey:PersonId" json:"-"`
	ConsumerServices        []ConsumerService        `gorm:"foreignKey:PersonId" json:"-"`
	Goals                   []Goal                   `gorm:"foreignKey:PersonId" json:"-"`
}
