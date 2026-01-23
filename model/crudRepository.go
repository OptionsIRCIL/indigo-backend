package model

import (
	"gorm.io/gorm"
)

//early mostly untested/unfinished crud commands

type CrudRepository[T any] struct {
	DB *gorm.DB
}

func (r *CrudRepository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *CrudRepository[T]) GetByID(id string) (*T, error) {
	var entity T
	if err := r.DB.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *CrudRepository[T]) GetAll() ([]T, error) {
	var results []T
	if err := r.DB.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *CrudRepository[T]) Update(id string, updates interface{}) error {
	var entity T
	return r.DB.Model(&entity).Where("id = ?", id).Updates(updates).Error
}

func (r *CrudRepository[T]) Delete(id string) error {
	var entity T
	return r.DB.Delete(&entity, "id = ?", id).Error
}
