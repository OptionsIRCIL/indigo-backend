package models

type Place struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	AddressLine1 string `gorm:"column:addressLine1;size:255;not null" json:"addressLine1"`
	AddressLine2 string `gorm:"column:addressLine2;size:255" json:"addressLine2"`
	City         string `gorm:"column:city;size:100;not null" json:"city"`
	State        string `gorm:"column:state;size:100;not null" json:"state"`
	ZipCode      string `gorm:"column:zipCode;size:20" json:"zipCode"`
	Country      string `gorm:"column:country;size:100;not null" json:"country"`
}
