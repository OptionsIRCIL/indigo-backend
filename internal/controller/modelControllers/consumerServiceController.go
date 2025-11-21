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

type ConsumerServiceController struct {
	Repo        *repository.BaseRepository[models.ConsumerService]
	SyncService *service.ExternalSyncService
}

func NewConsumerServiceController(repo *repository.BaseRepository[models.ConsumerService], syncService *service.ExternalSyncService) *ConsumerServiceController {
	return &ConsumerServiceController{
		Repo:        repo,
		SyncService: syncService,
	}
}

func (c *ConsumerServiceController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var cs models.ConsumerService
	if fetchErr := c.Repo.GetByID(uint(id), &cs); fetchErr != nil {
		log.Printf("Fetch error for ConsumerService %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncConsumerService(cs); syncErr != nil {
		log.Printf("External Sync Failed for ConsumerService %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("ConsumerService ID %s updated locally and synced externally.", idStr),
	})
}
