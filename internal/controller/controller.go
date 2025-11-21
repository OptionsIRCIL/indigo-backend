package controller

import (
	"myoptions.info/indigo/backend/internal/controller/modelControllers"
	"myoptions.info/indigo/backend/internal/repository"
	"myoptions.info/indigo/backend/internal/service"
)

// Controllers holds all initialized controller instances.
type Controllers struct {
	AddressPhone           *modelControllers.AddressPhoneController
	Alias                  *modelControllers.AliasController
	CommunityServiceEvent  *modelControllers.CommunityServiceEventController
	ConsumerService        *modelControllers.ConsumerServiceController
	DisabilityInfo         *modelControllers.DisabilityInfoController
	Goal                   *modelControllers.GoalController // We previously defined this one
	InformationAndReferral *modelControllers.InformationAndReferralController
	Organization           *modelControllers.OrganizationController
	Person                 *modelControllers.PersonController
	Place                  *modelControllers.PlaceController
	RecordDef              *modelControllers.RecordDefController
	ServicesOffered        *modelControllers.ServicesOfferedController
}

// NewControllers initializes all HTTP controllers, injecting their required
// dependencies (Repositories and Services).
func NewControllers(repos *repository.Repositories, services *service.Services) *Controllers {
	syncService := services.ExternalSync

	return &Controllers{
		AddressPhone:           modelControllers.NewAddressPhoneController(repos.AddressPhone, syncService),
		Alias:                  modelControllers.NewAliasController(repos.Alias, syncService),
		CommunityServiceEvent:  modelControllers.NewCommunityServiceEventController(repos.CommunityServiceEvent, syncService),
		ConsumerService:        modelControllers.NewConsumerServiceController(repos.ConsumerService, syncService),
		DisabilityInfo:         modelControllers.NewDisabilityInfoController(repos.DisabilityInfo, syncService),
		Goal:                   modelControllers.NewGoalController(repos.Goal, syncService),
		InformationAndReferral: modelControllers.NewInformationAndReferralController(repos.InformationAndReferral, syncService),
		Organization:           modelControllers.NewOrganizationController(repos.Organization, syncService),
		Person:                 modelControllers.NewPersonController(repos.Person, syncService),
		Place:                  modelControllers.NewPlaceController(repos.Place, syncService),
		RecordDef:              modelControllers.NewRecordDefController(repos.RecordDef, syncService),
		ServicesOffered:        modelControllers.NewServicesOfferedController(repos.ServicesOffered, syncService),
	}
}
