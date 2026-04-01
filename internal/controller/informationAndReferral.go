package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/model/entity"
)

func InformationAndReferralEffortPost(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		token := util.FetchTokenFromContext(r)
		employee := util.TokenToEmployee(token, database, ctx)

		// TODO: Users currently have the ability to create effort records for other employees if
		//       the employeeId key is specified. This should be restricted to only be allowed for administrators.
		deserializationErr, deserialized := util.Deserialize[entity.InformationAndReferralEffort](r.Body, []string{"post"})
		if deserializationErr != nil {
			// TODO: Are these messages safe to relay back to client?
			util.ThrowHttpError(w, 422, "Could not deserialize POST body: "+deserializationErr.Error())
			return
		}

		deserialized.InformationAndReferralId, _ = uuid.Parse(r.PathValue("informationAndReferralId"))

		// Set employee ID if not already provided
		if deserialized.EmployeeId.String() == "00000000-0000-0000-0000-000000000000" {
			deserialized.EmployeeId = employee.Id
		}

		createErr := gorm.G[entity.InformationAndReferralEffort](database).Create(ctx, &deserialized)
		if createErr != nil {
			util.ThrowHttpUnhandled(w, createErr)
			return
		}

		util.ReturnSerialized(w, 201, deserialized, []string{"get"})
	}
}
