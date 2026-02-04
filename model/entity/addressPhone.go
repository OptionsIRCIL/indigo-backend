package entity

type AddressPhone struct {
	Identifier
	//Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PhoneNumber string `gorm:"column:phoneNumber;size:25;not null" json:"phoneNumber"`
	Type        string `gorm:"column:type;size:50;not null" json:"type"`

	PlaceId  uint  `gorm:"column:placeId;not null" json:"placeId"`
	Place    Place `gorm:"foreignKey:PlaceId" json:"-"`
	PersonId uint  `gorm:"column:personId;not null" json:"personId"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:AddressPhoneId" json:"-"`
}
