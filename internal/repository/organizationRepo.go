package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type OrganizationRepository struct {
	*BaseRepository[models.Organization]
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	baseRepo := NewBaseRepository[models.Organization](db)
	return &OrganizationRepository{
		BaseRepository: baseRepo,
	}
}
