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

// AliasController provides the HTTP handlers for Alias operations.
type AliasController struct {
	Repo        *repository.AliasRepository
	SyncService *service.ExternalSyncService
}

func NewAliasController(repo *repository.AliasRepository, syncService *service.ExternalSyncService) *AliasController {
	return &AliasController{
		Repo:        repo,
		SyncService: syncService,
	}
}

// UpdateAndSyncHandler fetches the latest Alias record and pushes it externally.
func (c *AliasController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	// Fetch the data from the local DB
	var alias models.Alias
	if fetchErr := c.Repo.GetByID(uint(id), &alias); fetchErr != nil {
		log.Printf("Fetch error for Alias %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Alias record not found or database error.")
		return
	}

	// Trigger the External Sync
	if syncErr := c.SyncService.SyncAlias(alias); syncErr != nil {
		log.Printf("External Sync Failed for Alias %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	// Success Response
	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Alias ID %s updated locally and synced externally.", idStr),
	})
}
