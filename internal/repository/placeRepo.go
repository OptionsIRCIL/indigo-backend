package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type PlaceRepository struct {
	*BaseRepository[models.Place]
}

func NewPlaceRepository(db *gorm.DB) *PlaceRepository {
	baseRepo := NewBaseRepository[models.Place](db)
	return &PlaceRepository{
		BaseRepository: baseRepo,
	}
}
