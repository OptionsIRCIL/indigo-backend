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

type PersonController struct {
	Repo        *repository.BaseRepository[models.Person]
	SyncService *service.ExternalSyncService
}

func NewPersonController(repo *repository.BaseRepository[models.Person], syncService *service.ExternalSyncService) *PersonController {
	return &PersonController{
		Repo:        repo,
		SyncService: syncService,
	}
}

func (c *PersonController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var person models.Person
	if fetchErr := c.Repo.GetByID(uint(id), &person); fetchErr != nil {
		log.Printf("Fetch error for Person %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncPerson(person); syncErr != nil {
		log.Printf("External Sync Failed for Person %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Person ID %s updated locally and synced externally.", idStr),
	})
}
