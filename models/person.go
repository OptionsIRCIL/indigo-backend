package models

import "time"

type Person struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	FirstName  string    `gorm:"column:firstName;size:100;not null"`
	LastName   string    `gorm:"column:lastName;size:100;not null"`
	Salutation string    `gorm:"column:salutation;size:20"`
	Gender     string    `gorm:"column:gender;size:20"`
	Birthday   time.Time `gorm:"column:birthday;type:date"`
	Email      string    `gorm:"column:email;size:150;unique;not null"`
	Phone      string    `gorm:"column:phone;size:25"`
	Active     bool      `gorm:"column:active;default:true"`
	Deceased   bool      `gorm:"column:deceased;default:false"`

	AddressPhones           []AddressPhone           `gorm:"foreignKey:PersonID"`
	Aliases                 []Alias                  `gorm:"foreignKey:PersonID"`
	RecordDefs              []RecordDef              `gorm:"foreignKey:PersonID"`
	DisabilityInfo          *DisabilityInfo          `gorm:"foreignKey:PersonID"`
	InformationAndReferrals []InformationAndReferral `gorm:"foreignKey:PersonID"`
	ConsumerServices        []ConsumerService        `gorm:"foreignKey:PersonID"`
	Goals                   []Goal                   `gorm:"foreignKey:PersonID"`
}
