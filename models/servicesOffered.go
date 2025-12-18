package models

type ServicesOffered struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceName string `gorm:"column:serviceName;size:100;not null" json:"serviceName"`
	Description string `gorm:"column:description;size:255" json:"description"`

	ConsumerServiceID uint `gorm:"column:consumerServiceId;not null" json:"consumerServiceId"`
	AddressPhoneID    uint `gorm:"column:addressPhoneId;not null" json:"addressPhoneId"`
}
