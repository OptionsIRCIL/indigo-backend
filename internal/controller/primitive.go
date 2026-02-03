package controller

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/util"
)

// A PrimitiveFailure is returned when a primitive controller fails in either a non-fatal fashion.
type PrimitiveFailure struct {
	// A description of the encountered error
	Msg string
}

// SerializationParameters are used in the serialization/deserialization of database objects to
// mask properties with the `groups` tag.
type SerializationParameters struct {
	SerializationGroup   []string
	DeserializationGroup []string
}

// Error is implemented for conformity as an error.
func (p *PrimitiveFailure) Error() string {
	return p.Msg
}

// PrimitiveGetOne creates an http.HandlerFunc that GETs a single entity by ID.
// The URL for this handler expects one http.PathValue `id`.
func PrimitiveGetOne[Entity interface{}](database *gorm.DB, serializationGroups []string) http.HandlerFunc {
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

		// TODO: Conditional keys? Support via middleware?
		util.ReturnSerialized(w, 200, entity, serializationGroups)
	}
}

// PrimitiveGetCollection creates an http.HandlerFunc that GETs many entities based
// upon filter criteria stored in query parameters.
func PrimitiveGetCollection[E interface{}](database gorm.DB, serializationGroup string) http.HandlerFunc {
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
func PrimitivePost[Entity interface{}](database *gorm.DB, serializationParameters SerializationParameters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		deserializationGroups := serializationParameters.DeserializationGroup
		if deserializationGroups == nil {
			deserializationGroups = make([]string, 0)
		}

		deserializationErr, entity := util.Deserialize[Entity](r.Body, serializationParameters.DeserializationGroup)
		if deserializationErr != nil {
			util.ThrowHttpError(w, 422, "Could not deserialize POST body")
			return
		}

		createErr := gorm.G[Entity](database).Create(ctx, &entity)
		if createErr != nil {
			util.ThrowHttpUnhandled(w, createErr)
			return
		}

		util.ReturnSerialized(w, 201, entity, serializationParameters.SerializationGroup)
	}
}

// PrimitivePut acts similar to PrimitivePost, however, it overwrites
// an existing entity rather than creating a new one. The updated entity
// is echoed back in the response.
func PrimitivePut[E interface{}](database gorm.DB) http.HandlerFunc {
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
func PrimitiveDelete[Entity interface{}](database *gorm.DB) http.HandlerFunc {
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
