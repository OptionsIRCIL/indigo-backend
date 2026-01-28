package entity

import "time"

type Person struct {
	Identifier
	//Id         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName  string    `gorm:"column:firstName;size:100;not null" json:"firstName"`
	LastName   string    `gorm:"column:lastName;size:100;not null" json:"lastName"`
	Salutation string    `gorm:"column:salutation;size:20" json:"salutation"`
	Gender     string    `gorm:"column:gender;size:20" json:"gender"`
	Birthday   time.Time `gorm:"column:birthday;type:date" json:"birthday"`
	Email      string    `gorm:"column:email;size:150;unique;not null" json:"email"`
	Phone      string    `gorm:"column:phone;size:25" json:"phone"`
	Active     bool      `gorm:"column:active;default:true" json:"active"`
	Deceased   bool      `gorm:"column:deceased;default:false" json:"deceased"`

	AddressPhones           []AddressPhone           `gorm:"foreignKey:PersonId" json:"-"`
	Aliases                 []Alias                  `gorm:"foreignKey:PersonId" json:"-"`
	RecordDefs              []RecordDef              `gorm:"foreignKey:PersonId" json:"-"`
	DisabilityInfo          *DisabilityInfo          `gorm:"foreignKey:PersonId" json:"-"`
	InformationAndReferrals []InformationAndReferral `gorm:"foreignKey:PersonId" json:"-"`
	ConsumerServices        []ConsumerService        `gorm:"foreignKey:PersonId" json:"-"`
	Goals                   []Goal                   `gorm:"foreignKey:PersonId" json:"-"`
}
