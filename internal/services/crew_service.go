package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"iss-model-backend/internal/models"
)

const (
	CREW_URL        = "http://api.open-notify.org/astros.json"
	NASA_IMAGES_API = "https://images-api.nasa.gov/search"
)

var API_KEY = os.Getenv("NASA_API_KEY")

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

	var crewWithPhotos []models.AstronautWithPhoto

	for _, astronaut := range crew.People {
		photoURL, err := s.GetNASAAstronautPortrait(astronaut.Name)
		if err != nil {
			fmt.Errorf("failed to find correspondig image")
		}

		crewWithPhotos = append(crewWithPhotos, models.AstronautWithPhoto{
			Name:     astronaut.Name,
			ImageUrl: photoURL,
		})
	}

	return &models.ISSCrewWithPhotosResponse{
		People:  crewWithPhotos,
		Message: crew.Message,
	}, nil
}

func (s *CrewService) GetNASAAstronautPortrait(astronautName string) (string, error) {
	params := url.Values{}
	params.Set("q", astronautName+" portrait")
	params.Set("media_type", "image")
	params.Set("year_start", "2000")

	searchURL := NASA_IMAGES_API + "?" + params.Encode()

	resp, err := s.client.Get(searchURL)
	if err != nil {
		return "", fmt.Errorf("failed to search NASA images: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NASA Images API returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read NASA images response: %w", err)
	}

	var nasaResponse models.NASAImagesResponse
	if err := json.Unmarshal(body, &nasaResponse); err != nil {
		return "", fmt.Errorf("failed to parse NASA images response: %w", err)
	}

	return nasaResponse.Collection.Items[0].Links[0].Href, nil
}
