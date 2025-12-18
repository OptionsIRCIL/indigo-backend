package repository

import (
	"gorm.io/gorm"

	"myoptions.info/indigo/backend/models"
)

// Repositories holds all initialized repository instances.
type Repositories struct {
	AddressPhone           *BaseRepository[models.AddressPhone]
	Alias                  *BaseRepository[models.Alias]
	CommunityServiceEvent  *BaseRepository[models.CommunityServiceEvent]
	ConsumerService        *BaseRepository[models.ConsumerService]
	DisabilityInfo         *BaseRepository[models.DisabilityInfo]
	Goal                   *BaseRepository[models.Goal]
	InformationAndReferral *BaseRepository[models.InformationAndReferral]
	Organization           *BaseRepository[models.Organization]
	Person                 *BaseRepository[models.Person]
	Place                  *BaseRepository[models.Place]
	RecordDef              *BaseRepository[models.RecordDef]
	ServicesOffered        *BaseRepository[models.ServicesOffered]
}

// NewRepositories initializes all repositories using the provided GORM DB connection.
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		AddressPhone:           NewBaseRepository[models.AddressPhone](db),
		Alias:                  NewBaseRepository[models.Alias](db),
		CommunityServiceEvent:  NewBaseRepository[models.CommunityServiceEvent](db),
		ConsumerService:        NewBaseRepository[models.ConsumerService](db),
		DisabilityInfo:         NewBaseRepository[models.DisabilityInfo](db),
		Goal:                   NewBaseRepository[models.Goal](db),
		InformationAndReferral: NewBaseRepository[models.InformationAndReferral](db),
		Organization:           NewBaseRepository[models.Organization](db),
		Person:                 NewBaseRepository[models.Person](db),
		Place:                  NewBaseRepository[models.Place](db),
		RecordDef:              NewBaseRepository[models.RecordDef](db),
		ServicesOffered:        NewBaseRepository[models.ServicesOffered](db),
	}
}
