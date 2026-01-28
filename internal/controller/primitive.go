package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/util"
)

// A PrimitiveFailure is returned when a primitive controller fails in either a non-fatal fashion.
type PrimitiveFailure struct {
	// A description of the encountered error
	Msg string
}

// Error is implemented for conformity as an error.
func (p *PrimitiveFailure) Error() string {
	return p.Msg
}

func createUrlFormatExpr(urlFormat string) *regexp.Regexp {
	// TODO: Cap to EOL?
	// TODO: Support composite keys?
	// TODO: Expand to support multiple identifying columns, identifiers other than UUID
	// TODO: Unit tests
	idExpr, exprErr := regexp.Compile(
		"^" + strings.Replace(
			urlFormat,
			":id:",
			"([0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12})",
			1,
		),
	)

	if exprErr != nil {
		log.Fatal(exprErr)
	}

	return idExpr
}

// EXPERIMENTAL UUID REGEX HELPER
func extractId(r *http.Request, expr *regexp.Regexp) string {
	matches := expr.FindStringSubmatch(r.URL.Path)
	if len(matches) > 1 {
		return matches[1] // Return the captured UUID
	}
	return ""
}

//----- Primitive function section ------

// PrimitiveGetOne creates an http.HandlerFunc that GETs a single entity by ID.
// Format your urlFormat like this: `/path/to/my/:id:`, where `:id:` is an ID onto
// the desired entity.
func PrimitiveGetOne[Entity any](database *gorm.DB, urlFormat string) http.HandlerFunc {
	idExpr := createUrlFormatExpr(urlFormat)
	return func(w http.ResponseWriter, r *http.Request) {
		id := extractId(r, idExpr)
		if id == "" {
			util.ThrowHttpStatus(w, 404)
			return
		}

		// TODO: Fatal vs not-found
		var entity Entity
		if err := database.Where("id = ?", id).First(&entity).Error; err != nil {
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
func PrimitiveGetCollection[E any](database *gorm.DB) http.HandlerFunc {
	// TODO: Allowed properties for filtering?
	// TODO: Query parameter parsing
	// TODO: Database query
	return func(w http.ResponseWriter, r *http.Request) {
		var entities []E
		// Supports basic query filtering so we can expand later
		if err := database.Find(&entities).Error; err != nil {
			util.ThrowHttpUnhandled(w, err)
			return
		}
		util.ReturnSerialized(w, 200, entities)
	}
}

// PrimitivePost creates an http.HandlerFunc that accepts POST data containing
// a JSON-serialized entity and stores it to the database. The stored entity,
// including newly generated IDs, is echoed back in the response.
func PrimitivePost[E any](database *gorm.DB) http.HandlerFunc {
	// TODO: Validation (maybe method on struct?)
	return func(w http.ResponseWriter, r *http.Request) {
		var entity E
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			util.ThrowHttpStatus(w, 400)
			return
		}

		// GORM should automatically run the Identifier.BeforeCreate hook here
		if err := database.Create(&entity).Error; err != nil {
			util.ThrowHttpUnhandled(w, err)
			return
		}

		util.ReturnSerialized(w, 201, entity)
	}
}

// PrimitivePut acts similar to PrimitivePost, however, it overwrites
// an existing entity rather than creating a new one. The updated entity
// is echoed back in the response.
func PrimitivePut[E any](database *gorm.DB, urlFormat string) http.HandlerFunc {
	idExpr := createUrlFormatExpr(urlFormat)
	return func(w http.ResponseWriter, r *http.Request) {
		id := extractId(r, idExpr)
		if id == "" {
			util.ThrowHttpStatus(w, 404)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			util.ThrowHttpStatus(w, 400)
			return
		}

		// Update columns
		if err := database.Model(new(E)).Where("id = ?", id).Updates(updates).Error; err != nil {
			util.ThrowHttpUnhandled(w, err)
			return
		}

		w.WriteHeader(204)
	}
}

// PrimitiveDelete deletes the entity from the database given an ID stored in the URL.
// 204 is returned on successful deletion.
func PrimitiveDelete[Entity any](database *gorm.DB, urlFormat string) http.HandlerFunc {
	idExpr := createUrlFormatExpr(urlFormat)

	return func(w http.ResponseWriter, r *http.Request) {
		// Open context at start of function
		// TODO: Research context more to see if this is bad practice
		ctx := context.Background()

		// Pull ID
		id := idExpr.FindString(r.URL.Path)
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
