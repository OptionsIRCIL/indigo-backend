package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type ConsumerServiceRepository struct {
	*BaseRepository[models.ConsumerService]
}

func NewConsumerServiceRepository(db *gorm.DB) *ConsumerServiceRepository {
	baseRepo := NewBaseRepository[models.ConsumerService](db)
	return &ConsumerServiceRepository{
		BaseRepository: baseRepo,
	}
}

// add extra methods here
