package entity

type ConsumerService struct {
	Identifier
	//Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceName string `gorm:"column:serviceName;size:100;not null" json:"serviceName"`
	Status      string `gorm:"column:status;size:50" json:"status"`

	PersonId uint `gorm:"column:personId;not null" json:"personId"`

	ServicesOffered []ServicesOffered `gorm:"foreignKey:ConsumerServiceId" json:"-"`
}
