package models

type ServicesOffered struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	ServiceName string `gorm:"column:serviceName;size:100;not null"`
	Description string `gorm:"column:description;size:255"`

	ConsumerServiceID uint `gorm:"column:consumerServiceId;not null"`
	AddressPhoneID    uint `gorm:"column:addressPhoneId;not null"`
}
