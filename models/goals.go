package models

type Goal struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	GoalName    string `gorm:"column:goalName;size:100;not null" json:"goalName"`
	Description string `gorm:"column:description;size:255" json:"description"`
	Status      string `gorm:"column:status;size:50" json:"status"`

	PersonID uint `gorm:"column:personId;not null" json:"personId"`
}
