package entity

import (
	"time"
)

// A LocalUser is a user that uses local authentication (as opposed to LDAP auth)
// for application login.
type LocalUser struct {
	Username     string `gorm:"primaryKey;not null;size:255"`
	PasswordHash string `gorm:"not null;size:255"`
	ExpiresAt    *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
