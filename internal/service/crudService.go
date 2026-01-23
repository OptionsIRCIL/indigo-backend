package service

import "myoptions.info/indigo/backend/model"

//empty general crud service for creating special functions exclusive for individual models

type generalService[T any] struct {
	Repo *model.CrudRepository[T]
}

func (s *generalService[T]) Add(item *T) error {
	//add logic here
	return s.Repo.Create(item)
}
