package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type DisabilityInfoRepository struct {
	*BaseRepository[models.DisabilityInfo]
}

func NewDisabilityInfoRepository(db *gorm.DB) *DisabilityInfoRepository {
	baseRepo := NewBaseRepository[models.DisabilityInfo](db)
	return &DisabilityInfoRepository{
		BaseRepository: baseRepo,
	}
}

// extra methods here.
