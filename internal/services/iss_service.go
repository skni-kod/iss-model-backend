package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"iss-model-backend/internal/models"

	"gorm.io/gorm"
)

const (
	BASE_URL             = "https://api.wheretheiss.at/v1"
	ISS_ID               = 25544
	DATA_RETENTION_HOURS = 8
	COLLECTION_INTERVAL  = 10 * time.Second
	API_TIMEOUT          = 30 * time.Second
)

type ISSService struct {
	db *gorm.DB
}

func NewISSService(db *gorm.DB) *ISSService {
	service := &ISSService{db: db}

	go service.startDataCollection()

	go service.startCleanupRoutine()

	return service
}

func (s *ISSService) GetCurrentPosition(units string) (*models.ISSPosition, error) {
	if units == "" {
		units = "kilometers"
	}

	var recentPos models.ISSPosition
	cutoff := time.Now().Add(-30 * time.Second).Unix()

	err := s.db.Where("timestamp >= ?", cutoff).
		Order("timestamp desc").
		First(&recentPos).Error

	if err == nil {
		s.convertUnits(&recentPos, units)
		return &recentPos, nil
	}

	position, err := s.fetchFromAPI(0, units)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current position: %w", err)
	}

	s.db.Where("timestamp = ?", position.Timestamp).FirstOrCreate(position)

	return position, nil
}

func (s *ISSService) GetHistoricalPosition(timestamp int64, units string) (*models.ISSPosition, error) {
	if units == "" {
		units = "kilometers"
	}

	now := time.Now().Unix()
	if timestamp < now-4*3600 || timestamp > now+4*3600 {
		return nil, fmt.Errorf("timestamp outside retention window (4 hours back/forward)")
	}

	var position models.ISSPosition
	err := s.db.Raw("SELECT * FROM positions WHERE timestamp BETWEEN ? AND ? ORDER BY ABS(timestamp - ?) LIMIT 1",
		timestamp-60, timestamp+60, timestamp).
		Scan(&position).Error

	if err == nil {
		s.convertUnits(&position, units)
		return &position, nil
	}

	apiPosition, err := s.fetchFromAPI(timestamp, units)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical position: %w", err)
	}

	s.db.Where("timestamp = ?", apiPosition.Timestamp).FirstOrCreate(apiPosition)

	return apiPosition, nil
}

func (s *ISSService) GetPositionsInRange(startTime, endTime int64, units string) ([]*models.ISSPosition, error) {
	if units == "" {
		units = "kilometers"
	}

	var positions []*models.ISSPosition
	err := s.db.Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Order("timestamp asc").
		Find(&positions).Error
	if err != nil {
		return nil, err
	}

	for _, pos := range positions {
		s.convertUnits(pos, units)
	}

	return positions, nil
}

func (s *ISSService) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalCount int64
	if err := s.db.Model(&models.ISSPosition{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_positions"] = totalCount

	var latest models.ISSPosition
	if err := s.db.Order("timestamp desc").First(&latest).Error; err == nil {
		stats["latest_position"] = latest
		stats["latest_timestamp"] = time.Unix(latest.Timestamp, 0).Format(time.RFC3339)
	}

	var oldest models.ISSPosition
	if err := s.db.Order("timestamp asc").First(&oldest).Error; err == nil {
		stats["oldest_timestamp"] = time.Unix(oldest.Timestamp, 0).Format(time.RFC3339)
	}

	var recentCount int64
	oneHourAgo := time.Now().Add(-time.Hour).Unix()
	if err := s.db.Model(&models.ISSPosition{}).Where("timestamp >= ?", oneHourAgo).Count(&recentCount).Error; err == nil {
		stats["positions_last_hour"] = recentCount
	}

	if totalCount > 0 {
		stats["data_coverage"] = map[string]interface{}{
			"start": time.Unix(oldest.Timestamp, 0).Format(time.RFC3339),
			"end":   time.Unix(latest.Timestamp, 0).Format(time.RFC3339),
		}
	}

	return stats, nil
}

func (s *ISSService) fetchFromAPI(timestamp int64, units string) (*models.ISSPosition, error) {
	url := fmt.Sprintf("%s/satellites/%d", BASE_URL, ISS_ID)

	if timestamp > 0 {
		url += fmt.Sprintf("?timestamp=%d", timestamp)
		if units != "kilometers" {
			url += fmt.Sprintf("&units=%s", units)
		}
	} else if units != "kilometers" {
		url += fmt.Sprintf("?units=%s", units)
	}

	client := &http.Client{Timeout: API_TIMEOUT}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse models.ISSPositionResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	position := apiResponse.ToISSPosition()

	return position, nil
}

func (s *ISSService) convertUnits(position *models.ISSPosition, targetUnits string) {
	if position.Units == targetUnits {
		return
	}

	const KM_TO_MILES = 0.621371
	const MILES_TO_KM = 1.60934

	if position.Units == "kilometers" && targetUnits == "miles" {
		position.Altitude *= KM_TO_MILES
		position.Velocity *= KM_TO_MILES
		position.Footprint *= KM_TO_MILES
		position.Units = "miles"
	} else if position.Units == "miles" && targetUnits == "kilometers" {
		position.Altitude *= MILES_TO_KM
		position.Velocity *= MILES_TO_KM
		position.Footprint *= MILES_TO_KM
		position.Units = "kilometers"
	}
}

func (s *ISSService) startDataCollection() {
	log.Println("Starting ISS data collection...")

	ticker := time.NewTicker(COLLECTION_INTERVAL)
	defer ticker.Stop()

	s.collectData()

	for range ticker.C {
		s.collectData()
	}
}

func (s *ISSService) collectData() {
	position, err := s.fetchFromAPI(0, "kilometers")
	if err != nil {
		log.Printf("Failed to fetch ISS position: %v", err)
		return
	}

	var existingPos models.ISSPosition
	result := s.db.Where("timestamp = ?", position.Timestamp).FirstOrCreate(&existingPos, position)

	if result.Error != nil {
		log.Printf("Failed to store ISS position: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Stored new ISS position: lat=%.4f, lon=%.4f, timestamp=%d",
			position.Latitude, position.Longitude, position.Timestamp)
	}
}

func (s *ISSService) startCleanupRoutine() {
	log.Println("Starting ISS data cleanup routine...")

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupOldData()
	}
}

func (s *ISSService) cleanupOldData() {
	cutoff := time.Now().Add(-DATA_RETENTION_HOURS * time.Hour).Unix()

	result := s.db.Where("timestamp < ?", cutoff).Delete(&models.ISSPosition{})

	if result.Error != nil {
		log.Printf("Failed to cleanup old ISS positions: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d old ISS position records", result.RowsAffected)
	}
}
