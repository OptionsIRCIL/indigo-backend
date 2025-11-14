package repository

import (
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models"
)

type RecordDefRepository struct {
	*BaseRepository[models.RecordDef]
}

func NewRecordDefRepository(db *gorm.DB) *RecordDefRepository {
	baseRepo := NewBaseRepository[models.RecordDef](db)
	return &RecordDefRepository{
		BaseRepository: baseRepo,
	}
}
