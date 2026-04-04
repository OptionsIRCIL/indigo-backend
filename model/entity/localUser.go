package entity

import (
	"time"
)

// A LocalUser is a user that uses local authentication (as opposed to LDAP auth)
// for application login.
type LocalUser struct {
	Username     string `gorm:"primaryKey; type: VARCHAR(255); not null"`
	PasswordHash string `gorm:"type: VARCHAR(255); not null"`
	ExpiresAt    *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
