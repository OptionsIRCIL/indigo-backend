package entity

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id        uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string    `gorm:"size:150;not null" json:"name"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
