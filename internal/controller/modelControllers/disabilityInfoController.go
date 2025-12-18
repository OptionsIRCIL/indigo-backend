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

type DisabilityInfoController struct {
	Repo        *repository.BaseRepository[models.DisabilityInfo]
	SyncService *service.ExternalSyncService
}

func NewDisabilityInfoController(repo *repository.BaseRepository[models.DisabilityInfo], syncService *service.ExternalSyncService) *DisabilityInfoController {
	return &DisabilityInfoController{
		Repo:        repo,
		SyncService: syncService,
	}
}

func (c *DisabilityInfoController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var di models.DisabilityInfo
	if fetchErr := c.Repo.GetByID(uint(id), &di); fetchErr != nil {
		log.Printf("Fetch error for DisabilityInfo %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncDisabilityInfo(di); syncErr != nil {
		log.Printf("External Sync Failed for DisabilityInfo %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("DisabilityInfo ID %s updated locally and synced externally.", idStr),
	})
}
