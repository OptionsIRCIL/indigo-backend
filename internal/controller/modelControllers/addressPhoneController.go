package modelControllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"myoptions.info/indigo/backend/internal/repository"
	"myoptions.info/indigo/backend/internal/service"
	u "myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/models"
)

// AddressPhoneController provides the HTTP handlers for AddressPhone operations.
type AddressPhoneController struct {
	Repo        *repository.AddressPhoneRepository
	SyncService *service.ExternalSyncService
}

func NewAddressPhoneController(repo *repository.AddressPhoneRepository, syncService *service.ExternalSyncService) *AddressPhoneController {
	return &AddressPhoneController{
		Repo:        repo,
		SyncService: syncService,
	}
}

// UpdateAndSyncHandler fetches the latest AddressPhone record and pushes it externally.
func (c *AddressPhoneController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	// Assuming the ID is passed via a query parameter for simplicity in this example
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	// Fetch the data from the local DB (Repository Layer)
	var ap models.AddressPhone
	if fetchErr := c.Repo.GetByID(uint(id), &ap); fetchErr != nil {
		log.Printf("Fetch error for AddressPhone %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	// Trigger the External Sync (Service Layer)
	if syncErr := c.SyncService.SyncAddressPhone(ap); syncErr != nil {
		log.Printf("External Sync Failed for AddressPhone %d: %v", id, syncErr)
		// Return a generic error to the client, but log the specific failure
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	// Success Response
	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("AddressPhone ID %s updated locally and synced externally.", idStr),
	})
}
