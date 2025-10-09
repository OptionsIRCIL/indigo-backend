package models

type ConsumerService struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	ServiceName string `gorm:"column:serviceName;size:100;not null"`
	Status      string `gorm:"column:status;size:50"`

	PersonID uint `gorm:"column:personId;not null"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:ConsumerServiceID"`
}
