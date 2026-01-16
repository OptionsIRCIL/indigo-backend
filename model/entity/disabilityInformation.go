package entity

type DisabilityInfo struct {
	Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Disability  string `gorm:"column:disability;size:100;not null" json:"disability"`
	Description string `gorm:"column:description;size:255" json:"description"`
	Severity    string `gorm:"column:severity;size:50" json:"severity"`

	PersonId uint `gorm:"column:personId;unique;not null" json:"personId"`
}
