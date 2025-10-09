package models

import "time"

type Organization struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;size:150;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`
}
