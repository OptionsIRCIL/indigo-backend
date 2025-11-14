package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type AddressPhoneRepository struct {
	*BaseRepository[models.AddressPhone]
}

// NewAddressPhoneRepository initializes the repository
func NewAddressPhoneRepository(db *gorm.DB) *AddressPhoneRepository {
	baseRepo := NewBaseRepository[models.AddressPhone](db)
	return &AddressPhoneRepository{
		BaseRepository: baseRepo,
	}
}

// add methods here
