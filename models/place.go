package models

type Place struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	AddressLine1 string `gorm:"column:addressLine1;size:255;not null"`
	AddressLine2 string `gorm:"column:addressLine2;size:255"`
	City         string `gorm:"column:city;size:100;not null"`
	State        string `gorm:"column:state;size:100;not null"`
	ZipCode      string `gorm:"column:zipCode;size:20"`
	Country      string `gorm:"column:country;size:100;not null"`
}
