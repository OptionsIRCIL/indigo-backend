package entity

import "time"

type CommunityServiceEvent struct {
	Id                      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InitialDate             time.Time `gorm:"column:initialDate;type:datetime;not null" json:"initialDate"`
	Category                string    `gorm:"column:category;size:100;not null" json:"category"`
	FutureReference         string    `gorm:"column:futureReference;size:255" json:"futureReference"`
	ServiceDescription      string    `gorm:"column:serviceDescription;size:255;not null" json:"serviceDescription"`
	Outcome                 string    `gorm:"column:outcome;size:255" json:"outcome"`
	Publications            int       `gorm:"column:publications;default:0" json:"publications"`
	PersonsWithDisabilities int       `gorm:"column:personsWithDisabilities;default:0" json:"personsWithDisabilities"`
	GeneralPublic           int       `gorm:"column:generalPublic;default:0" json:"generalPublic"`
}
