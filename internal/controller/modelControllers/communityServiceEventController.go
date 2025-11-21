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

type CommunityServiceEventController struct {
	Repo        *repository.BaseRepository[models.CommunityServiceEvent]
	SyncService *service.ExternalSyncService
}

func NewCommunityServiceEventController(repo *repository.BaseRepository[models.CommunityServiceEvent], syncService *service.ExternalSyncService) *CommunityServiceEventController {
	return &CommunityServiceEventController{
		Repo:        repo,
		SyncService: syncService,
	}
}
func (c *CommunityServiceEventController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var event models.CommunityServiceEvent
	if fetchErr := c.Repo.GetByID(uint(id), &event); fetchErr != nil {
		log.Printf("Fetch error for CommunityServiceEvent %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncCommunityServiceEvent(event); syncErr != nil {
		log.Printf("External Sync Failed for CommunityServiceEvent %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("CommunityServiceEvent ID %s updated locally and synced externally.", idStr),
	})
}
