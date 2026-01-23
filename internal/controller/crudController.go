package controller

import (
	"encoding/json"
	"net/http"

	"myoptions.info/indigo/backend/model"
)

type CrudController[T any] struct {
	Repo *model.CrudRepository[T]
}

func (gc *CrudController[T]) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entity T
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, "Invalid Body", 400)
			return
		}
		gc.Repo.Create(&entity)
		json.NewEncoder(w).Encode(entity)
	}
}

func (gc *CrudController[T]) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") // Go 1.22+ feature
		res, err := gc.Repo.GetByID(id)
		if err != nil {
			http.Error(w, "Not Found", 404)
			return
		}
		json.NewEncoder(w).Encode(res)
	}
}

func (gc *CrudController[T]) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := gc.Repo.GetAll()
		json.NewEncoder(w).Encode(res)
	}
}

func (gc *CrudController[T]) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var up map[string]any
		json.NewDecoder(r.Body).Decode(&up)
		gc.Repo.Update(id, up)
		w.WriteHeader(204)
	}
}

func (gc *CrudController[T]) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gc.Repo.Delete(r.PathValue("id"))
		w.WriteHeader(204)
	}
}
