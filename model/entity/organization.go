package entity

import "time"

type Organization struct {
	Identifier
	//Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:150;not null" json:"name"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime" json:"updatedAt"`
}
