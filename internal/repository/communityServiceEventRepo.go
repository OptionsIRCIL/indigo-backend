package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type CommunityServiceEventRepository struct {
	*BaseRepository[models.CommunityServiceEvent]
}

// NewCommunityServiceEventRepository initializes the repository
func NewCommunityServiceEventRepository(db *gorm.DB) *CommunityServiceEventRepository {
	baseRepo := NewBaseRepository[models.CommunityServiceEvent](db)
	return &CommunityServiceEventRepository{
		BaseRepository: baseRepo,
	}
}

// add methods here
