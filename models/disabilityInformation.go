package models

type DisabilityInfo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Disability  string `gorm:"column:disability;size:100;not null"`
	Description string `gorm:"column:description;size:255"`
	Severity    string `gorm:"column:severity;size:50"`

	PersonID uint `gorm:"column:personId;unique;not null"`
}
