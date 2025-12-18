package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type InformationAndReferralRepository struct {
	*BaseRepository[models.InformationAndReferral]
}

func NewInformationAndReferralRepository(db *gorm.DB) *InformationAndReferralRepository {
	baseRepo := NewBaseRepository[models.InformationAndReferral](db)
	return &InformationAndReferralRepository{
		BaseRepository: baseRepo,
	}
}
