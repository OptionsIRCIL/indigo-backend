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

type GoalController struct {
	Repo        *repository.BaseRepository[models.Goal]
	SyncService *service.ExternalSyncService
}

func NewGoalController(repo *repository.BaseRepository[models.Goal], syncService *service.ExternalSyncService) *GoalController {
	return &GoalController{
		Repo:        repo,
		SyncService: syncService,
	}
}

func (c *GoalController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var goal models.Goal
	if fetchErr := c.Repo.GetByID(uint(id), &goal); fetchErr != nil {
		log.Printf("Fetch error for Goal %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncGoal(goal); syncErr != nil {
		log.Printf("External Sync Failed for Goal %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Goal ID %s updated locally and synced externally.", idStr),
	})
}
