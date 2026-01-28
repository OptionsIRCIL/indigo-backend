package entity

type ServicesOffered struct {
	Identifier
	//Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceName string `gorm:"column:serviceName;size:100;not null" json:"serviceName"`
	Description string `gorm:"column:description;size:255" json:"description"`

	ConsumerServiceId uint `gorm:"column:consumerServiceId;not null" json:"consumerServiceId"`
	AddressPhoneId    uint `gorm:"column:addressPhoneId;not null" json:"addressPhoneId"`
}
