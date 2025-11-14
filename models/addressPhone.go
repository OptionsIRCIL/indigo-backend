package models

type AddressPhone struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PhoneNumber string `gorm:"column:phoneNumber;size:25;not null" json:"phoneNumber"`
	Type        string `gorm:"column:type;size:50;not null" json:"type"`

	PlaceID  uint  `gorm:"column:placeId;not null" json:"placeId"`
	Place    Place `gorm:"foreignKey:PlaceID" json:"-"`
	PersonID uint  `gorm:"column:personId;not null" json:"personId"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:AddressPhoneID" json:"-"`
}
