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

type InformationAndReferralController struct {
	Repo        *repository.InfoAndReferralRepository
	SyncService *service.ExternalSyncService
}

func NewInformationAndReferralController(repo *repository.InformationAndReferralRepository, syncService *service.ExternalSyncService) *InformationAndReferralController {
	return &InformationAndReferralController{
		Repo:        repo,
		SyncService: syncService,
	}
}

func (c *InformationAndReferralController) UpdateAndSyncHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		u.ThrowHttpError(w, http.StatusBadRequest, "Invalid ID format or missing query parameter 'id'")
		return
	}

	var ir models.InformationAndReferral
	if fetchErr := c.Repo.GetByID(uint(id), &ir); fetchErr != nil {
		log.Printf("Fetch error for InformationAndReferral %d: %v", id, fetchErr)
		u.ThrowHttpError(w, http.StatusNotFound, "Record not found or database error.")
		return
	}

	if syncErr := c.SyncService.SyncInformationAndReferral(ir); syncErr != nil {
		log.Printf("External Sync Failed for InformationAndReferral %d: %v", id, syncErr)
		u.ThrowHttpError(w, http.StatusInternalServerError, "Request completed, but external data synchronization failed.")
		return
	}

	u.ReturnSerialized(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("InformationAndReferral ID %s updated locally and synced externally.", idStr),
	})
}
