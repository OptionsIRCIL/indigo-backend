package service

import "log"

// Services holds all initialized services
type Services struct {
	ExternalSync *ExternalSyncService
}

// NewServices initializes all services.
func NewServices(baseURL, authToken string) *Services {
	// ExternalSyncService will need configuration for full implementation
	log.Println("Initializing External Sync Service...")

	externalSyncService := NewExternalSyncService(baseURL, authToken)
	// ExternalSyncService is initialized
	return &Services{
		ExternalSync: externalSyncService,
	}
}
