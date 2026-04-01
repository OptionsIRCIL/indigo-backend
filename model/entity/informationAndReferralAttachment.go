package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InformationAndReferralAttachment struct {
	Id       uuid.UUID `gorm:"primaryKey; type: CHAR(36)"`
	FileName string    `gorm:"type: VARCHAR(255); not null"`

	// The file's uploader
	EmployeeId uuid.UUID `gorm:"not null"`
	Employee   Employee

	// MIME type for use in HEAD operations
	ContentType string `gorm:"type: VARCHAR(255); not null"`

	// Content size in bytes for HEAD
	Size uint `gorm:"type: BIGINT UNSIGNED; not null"`

	// SHA512 signature of the file
	Signature []byte `gorm:"type: BINARY(64); not null"`

	CreatedAt time.Time `gorm:"not null"`
	DeletedAt time.Time
}

func (i *InformationAndReferralAttachment) BeforeCreate(tx *gorm.DB) (err error) {
	i.Id, err = uuid.NewRandom()
	return err
}
