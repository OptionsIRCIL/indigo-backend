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

type PlaceController struct {
	Repo        *repository.PlaceRepository
	SyncService *service.ExternalSyncService
}

func NewPlaceController(repo *repository.PlaceRepository, syncService *service.ExternalSyncService) *PlaceController {
	return &PlaceController{
		Repo:        repo,
		SyncService: syncService,
	}
}

func (c *PlaceController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var p models.Place
	if fetchErr := c.Repo.GetByID(uint(id), &p); fetchErr != nil {
		log.Printf("Fetch error for Place %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncPlace(p); syncErr != nil {
		log.Printf("External Sync Failed for Place %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Place ID %s updated locally and synced externally.", idStr),
	})
}
