package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// BaseModel sets required vars for repos
type BaseModel interface {
	getID() uint
}

// BaseRepository defines the generic CRUD operations for all models.
// T is assigned as the pointer representing the models
type BaseRepository[T any] struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new generic repository instance.
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		DB: db,
	}
}

// Create inserts a new record into the database.
func (r *BaseRepository[T]) Create(entity *T) error {
	result := r.DB.Create(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to create record: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a single record by its primary key (ID).
// entity must be pointer for a struct/model.
func (r *BaseRepository[T]) GetByID(id uint, entity *T) error {
	result := r.DB.First(entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("record with ID %d not found", id)
		}
		return fmt.Errorf("failed to retrieve record by ID %d: %w", id, result.Error)
	}
	return nil
}

// Update saves changes to an existing record.
func (r *BaseRepository[T]) Update(entity *T) error {
	// Note: Save updates all fields.
	result := r.DB.Save(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update record: %w", result.Error)
	}
	return nil
}

// Delete removes a record based on its ID.
func (r *BaseRepository[T]) Delete(id uint) error {
	// Create an empty instance of the model type for deletion.
	var entity T
	result := r.DB.Delete(&entity, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete record by ID %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record with ID %d not found for deletion", id)
	}
	return nil
}
