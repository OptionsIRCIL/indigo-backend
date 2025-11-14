package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type AliasRepository struct {
	*BaseRepository[models.Alias]
}

// NewAliasRepository initializes the repository
func NewAliasRepository(db *gorm.DB) *AliasRepository {
	baseRepo := NewBaseRepository[models.Alias](db)
	return &AliasRepository{
		BaseRepository: baseRepo,
	}
}

// add methods here
