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

type OrganizationController struct {
	Repo        *repository.OrganizationRepository
	SyncService *service.ExternalSyncService
}

func NewOrganizationController(repo *repository.OrganizationRepository, syncService *service.ExternalSyncService) *OrganizationController {
	return &OrganizationController{
		Repo:        repo,
		SyncService: syncService,
	}
}

func (c *OrganizationController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var org models.Organization
	if fetchErr := c.Repo.GetByID(uint(id), &org); fetchErr != nil {
		log.Printf("Fetch error for Organization %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncOrganization(org); syncErr != nil {
		log.Printf("External Sync Failed for Organization %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Organization ID %s updated locally and synced externally.", idStr),
	})
}
