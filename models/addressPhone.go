package models

type AddressPhone struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	PhoneNumber string `gorm:"column:phoneNumber;size:25;not null"`
	Type        string `gorm:"column:type;size:50;not null"`

	PlaceID  uint  `gorm:"column:placeId;not null"`
	Place    Place `gorm:"foreignKey:PlaceID"`
	PersonID uint  `gorm:"column:personId;not null"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:AddressPhoneID"`
}
