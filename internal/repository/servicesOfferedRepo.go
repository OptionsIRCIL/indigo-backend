package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type ServicesOfferedRepository struct {
	*BaseRepository[models.ServicesOffered]
}

func NewServicesOfferedRepository(db *gorm.DB) *ServicesOfferedRepository {
	baseRepo := NewBaseRepository[models.ServicesOffered](db)
	return &ServicesOfferedRepository{
		BaseRepository: baseRepo,
	}
}
