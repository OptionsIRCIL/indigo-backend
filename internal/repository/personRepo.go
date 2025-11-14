package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type PersonRepository struct {
	*BaseRepository[models.Person]
}

func NewPersonRepository(db *gorm.DB) *PersonRepository {
	baseRepo := NewBaseRepository[models.Person](db)
	return &PersonRepository{
		BaseRepository: baseRepo,
	}
}
