package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"myoptions.info/indigo/backend/models"
)

// ExternalSyncService handles all outbound communication to the external API.
type ExternalSyncService struct {
	Client       *http.Client
	ApiBaseUrl   string
	ApiAuthToken string
}

// NewExternalSyncService initializes the service with an HTTP client.
func NewExternalSyncService(baseURL, authToken string) *ExternalSyncService {
	return &ExternalSyncService{
		Client:       &http.Client{Timeout: 10 * time.Second},
		ApiBaseUrl:   baseURL,
		ApiAuthToken: authToken,
	}
}

// pushDataToExternalAPI handles the common HTTP request, headers, and error checking.
func (s *ExternalSyncService) pushDataToExternalAPI(jsonData []byte, endpoint string) error {
	url := s.ApiBaseUrl + endpoint
	reqBody := bytes.NewReader(jsonData)

	req, err := http.NewRequest(http.MethodPut, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", endpoint, err)
	}

	// Set headers (Authentication and Content Type)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.ApiAuthToken)

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	// Check for status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("external API sync failed for %s. Status: %d, Response: %s",
			endpoint, resp.StatusCode, string(errorBody))
	}

	return nil
}

// SyncAddressPhone serializes the AddressPhone struct and pushes it to the external API.
func (s *ExternalSyncService) SyncAddressPhone(entity models.AddressPhone) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize AddressPhone: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/address-phones/%d", entity.ID) // Example PUT endpoint
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncAlias serializes the Alias struct and pushes it to the external API.
func (s *ExternalSyncService) SyncAlias(entity models.Alias) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize Alias: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/aliases/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncCommunityServiceEvent serializes the CommunityServiceEvent struct and pushes it to the external API.
func (s *ExternalSyncService) SyncCommunityServiceEvent(entity models.CommunityServiceEvent) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize CommunityServiceEvent: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/community-events/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncConsumerService serializes the ConsumerService struct and pushes it to the external API.
func (s *ExternalSyncService) SyncConsumerService(entity models.ConsumerService) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize ConsumerService: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/consumer-services/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncDisabilityInfo serializes the DisabilityInfo struct and pushes it to the external API.
func (s *ExternalSyncService) SyncDisabilityInfo(entity models.DisabilityInfo) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize DisabilityInfo: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/disability-info/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncGoal serializes the Goal struct and pushes it to the external API.
func (s *ExternalSyncService) SyncGoal(entity models.Goal) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize Goal: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/goals/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncInformationAndReferral serializes the InformationAndReferral struct and pushes it to the external API.
func (s *ExternalSyncService) SyncInformationAndReferral(entity models.InformationAndReferral) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize InformationAndReferral: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/info-referrals/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncOrganization serializes the Organization struct and pushes it to the external API.
func (s *ExternalSyncService) SyncOrganization(entity models.Organization) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize Organization: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/organizations/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncPerson serializes the Person struct and pushes it to the external API.
func (s *ExternalSyncService) SyncPerson(entity models.Person) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize Person: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/persons/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncPlace serializes the Place struct and pushes it to the external API.
func (s *ExternalSyncService) SyncPlace(entity models.Place) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize Place: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/places/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncRecordDef serializes the RecordDef struct and pushes it to the external API.
func (s *ExternalSyncService) SyncRecordDef(entity models.RecordDef) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize RecordDef: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/record-defs/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}

// SyncServicesOffered serializes the ServicesOffered struct and pushes it to the external API.
func (s *ExternalSyncService) SyncServicesOffered(entity models.ServicesOffered) error {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to serialize ServicesOffered: %w", err)
	}
	endpoint := fmt.Sprintf("/v1/services-offered/%d", entity.ID)
	return s.pushDataToExternalAPI(jsonData, endpoint)
}
