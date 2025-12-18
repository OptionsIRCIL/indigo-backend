package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type GoalRepository struct {
	*BaseRepository[models.Goal]
}

func NewGoalRepository(db *gorm.DB) *GoalRepository {
	baseRepo := NewBaseRepository[models.Goal](db)
	return &GoalRepository{
		BaseRepository: baseRepo,
	}
}
