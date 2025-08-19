package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"iss-model-backend/internal/models"
)

const (
	CREW_URL = "http://api.open-notify.org/astros.json"
)

type CrewService struct {
	client *http.Client
}

func NewCrewService() *CrewService {
	return &CrewService{
		client: &http.Client{Timeout: API_TIMEOUT},
	}
}

func (s *CrewService) GetCurrentCrew() (*models.ISSCrewResponse, error) {
	resp, err := s.client.Get(CREW_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current ISS crew: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch API: %s - %s", resp.Status, string(body))
	}

	var apiResponse models.ISSCrewResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response %w", err)
	}

	var issCrewMembers []models.Astronaut
	for _, person := range apiResponse.People {
		if person.Craft == "ISS" {
			issCrewMembers = append(issCrewMembers, person)
		}
	}

	apiResponse.People = issCrewMembers

	return &apiResponse, nil
}

func (s *CrewService) GetCurrentCrewWithPhotos() (*models.ISSCrewWithPhotosResponse, error) {
	crew, err := s.GetCurrentCrew()
	if err != nil {
		return nil, fmt.Errorf("failed to get crew: %w", err)
	}

	params := url.Values{}
	params.Set("action", "query")
	params.Set("prop", "images")
	params.Set("titles", crew.People[0].Name)
	params.Set("format", "json")

	resp, err := s.client.Get("https://en.wikipedia.org/w/api.php" + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current ISS crew: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf(string(body))

	return nil, nil
}
