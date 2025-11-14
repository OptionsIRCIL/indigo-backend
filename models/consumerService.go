package models

type ConsumerService struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceName string `gorm:"column:serviceName;size:100;not null" json:"serviceName"`
	Status      string `gorm:"column:status;size:50" json:"status"`

	PersonID uint `gorm:"column:personId;not null" json:"personId"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:ConsumerServiceID" json:"-"`
}
