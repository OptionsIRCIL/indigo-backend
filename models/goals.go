package models

type Goal struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	GoalName    string `gorm:"column:goalName;size:100;not null"`
	Description string `gorm:"column:description;size:255"`
	Status      string `gorm:"column:status;size:50"`

	PersonID uint `gorm:"column:personId;not null"`
}
