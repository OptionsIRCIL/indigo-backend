package models

import "time"

type CommunityServiceEvent struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement"`
	InitialDate             time.Time `gorm:"column:initialDate;type:datetime;not null"`
	Category                string    `gorm:"column:category;size:100;not null"`
	FutureReference         string    `gorm:"column:futureReference;size:255"`
	ServiceDescription      string    `gorm:"column:serviceDescription;size:255;not null"`
	Outcome                 string    `gorm:"column:outcome;size:255"`
	Publications            int       `gorm:"column:publications;default:0"`
	PersonsWithDisabilities int       `gorm:"column:personsWithDisabilities;default:0"`
	GeneralPublic           int       `gorm:"column:generalPublic;default:0"`
}
