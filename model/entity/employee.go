package entity

import "github.com/google/uuid"

type Employee struct {
	Id        uuid.UUID `gorm:"primaryKey;type:char(36)" json:"id"`
	FirstName string    `gorm:"size:255;" json:"firstName"`
	LastName  string    `gorm:"size:255;" json:"lastName"`

	// optional/may remove

	Username string `gorm:"size:255;" json:"username"`
	Email    string `gorm:"size:255;" json:"email"`
}
