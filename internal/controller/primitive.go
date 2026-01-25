package controller

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/util"
)

// A Primitive is one that is registered with GORM and has basic properties available for control.
type Primitive struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// A PrimitiveFailure is returned when a primitive controller fails in either a non-fatal fashion.
type PrimitiveFailure struct {
	// A description of the encountered error
	Msg string
}

// Error is implemented for conformity as an error.
func (p *PrimitiveFailure) Error() string {
	return p.Msg
}

// PrimitiveGetOne creates an http.HandlerFunc that GETs a single entity by ID.
func PrimitiveGetOne[Entity Primitive](database *gorm.DB, urlFormat string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open context at start of function
		// TODO: Research context more to see if this is bad practice
		ctx := context.Background()

		// Pull ID
		id := r.PathValue("id")
		if id == "" {
			util.ThrowHttpStatus(w, 404)
			return
		}

		// Lookup
		// This assumes that the id column is unique
		entity, err := gorm.G[Entity](database).Where("id = ?", id).First(ctx)

		// TODO: Fatal vs not-found
		if err != nil {
			util.ThrowHttpStatus(w, 404)
			return
		}

		// TODO: Support key selection?
		// TODO: Conditional keys? Support via middleware?
		util.ReturnSerialized(w, 200, entity)
	}
}

// PrimitiveGetCollection creates an http.HandlerFunc that GETs many entities based
// upon filter criteria stored in query parameters.
func PrimitiveGetCollection[E Primitive](database gorm.DB) http.HandlerFunc {
	// TODO: Allowed properties for filtering?
	// TODO: Query parameter parsing
	// TODO: Database query
	return func(w http.ResponseWriter, r *http.Request) {
		util.ThrowHttpUnhandled(
			w,
			&PrimitiveFailure{Msg: "GET collection method not yet implemented"},
		)
	}
}

// PrimitivePost creates an http.HandlerFunc that accepts POST data containing
// a JSON-serialized entity and stores it to the database. The stored entity,
// including newly generated IDs, is echoed back in the response.
func PrimitivePost[Entity Primitive](database *gorm.DB, group string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		deserializationErr, entity := util.Deserialize[Entity](r.Body, group)
		if deserializationErr != nil {
			util.ThrowHttpError(w, 422, "Could not deserialize POST body")
			return
		}

		// TODO ID generation
		createErr := gorm.G[Entity](database).Create(ctx, &entity)
		if createErr != nil {
			util.ThrowHttpUnhandled(w, createErr)
			return
		}

		util.ReturnSerialized(w, 201, entity)

	}
}

// PrimitivePut acts similar to PrimitivePost, however, it overwrites
// an existing entity rather than creating a new one. The updated entity
// is echoed back in the response.
func PrimitivePut[E Primitive](database gorm.DB) http.HandlerFunc {
	// TODO: Implement
	return func(w http.ResponseWriter, r *http.Request) {
		util.ThrowHttpUnhandled(
			w,
			&PrimitiveFailure{Msg: "PUT method not yet implemented"},
		)
	}
}

// PrimitiveDelete deletes the entity from the database given an ID stored in the URL.
// 204 is returned on successful deletion.
func PrimitiveDelete[Entity Primitive](database *gorm.DB, urlFormat string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open context at start of function
		// TODO: Research context more to see if this is bad practice
		ctx := context.Background()

		// Pull ID
		id := r.PathValue("id")
		if id == "" {
			util.ThrowHttpStatus(w, 404)
			return
		}

		// Begin transaction
		err := database.Transaction(func(tx *gorm.DB) error {
			// This assumes that the id column is unique
			_, err := gorm.G[Entity](tx).Where("id = ?", id).First(ctx)

			// TODO: Fatal vs not-found
			if err != nil {
				util.ThrowHttpStatus(w, 404)
				return err
			}

			// TODO: Soft deletion. Common interface?
			// https://gorm.io/docs/delete.html#Soft-Delete?
			rowsAffected, err := gorm.G[Entity](tx).Where("id = ?", id).Delete(ctx)

			if err != nil {
				return err
			}

			// Ensure only one row deleted
			if rowsAffected != 1 {
				return &PrimitiveFailure{Msg: "Multiple rows affected on delete operation, rolling back"}
			}

			// I think we're good chat
			return nil
		})

		// TODO: Better error messages/handling
		if err != nil {
			util.ThrowHttpUnhandled(w, err)
			return
		}

		// Operation was successful, return no content
		w.WriteHeader(204)
	}
}
