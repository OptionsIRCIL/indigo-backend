package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonAttachment struct {
	Id       uuid.UUID `gorm:"primaryKey; type: CHAR(36)" groups:"get"`
	FileName string    `gorm:"type: VARCHAR(255); not null" groups:"get"`

	// The file's uploader
	EmployeeId uuid.UUID `gorm:"not null" groups:"get"`
	Employee   Employee

	// The parent person
	PersonId uuid.UUID `gorm:"not null" groups:"get"`
	Person   Person

	// MIME type for use in HEAD operations
	ContentType string `gorm:"type: VARCHAR(255); not null" groups:"get"`

	// Content size in bytes for HEAD
	Size uint `gorm:"type: BIGINT UNSIGNED; not null" groups:"get"`

	// SHA512 signature of the file
	Signature string `gorm:"type: BINARY(64); not null"`

	CreatedAt time.Time `gorm:"not null" groups:"get"`
	DeletedAt time.Time
}

func (i *PersonAttachment) BeforeCreate(tx *gorm.DB) (err error) {
	if &i.Id == nil {
		i.Id, err = uuid.NewRandom()
		return err
	}
	return nil
}
