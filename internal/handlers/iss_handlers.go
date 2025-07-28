// internal/handlers/iss_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"iss-model-backend/internal/models"
	"iss-model-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type ISSHandler struct {
	issService *services.ISSService
}

func NewISSHandler(issService *services.ISSService) *ISSHandler {
	return &ISSHandler{
		issService: issService,
	}
}

// GetCurrentPosition returns the current ISS position
// @Summary Get Current ISS Position
// @Description Returns the current position of the International Space Station
// @Tags ISS
// @Accept json
// @Produce json
// @Param units query string false "Units (kilometers or miles)" Enums(kilometers, miles) default(kilometers)
// @Success 200 {object} models.ISSPosition
// @Failure 500 {object} models.ErrorResponse
// @Router /iss/current [get]
func (h *ISSHandler) GetCurrentPosition(w http.ResponseWriter, r *http.Request) {
	units := r.URL.Query().Get("units")
	if units != "miles" {
		units = "kilometers"
	}

	position, err := h.issService.GetCurrentPosition(units)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get current position", err.Error())
		return
	}

	h.sendJSONResponse(w, http.StatusOK, position)
}

// GetHistoricalPosition returns ISS position for a specific timestamp
// @Summary Get Historical ISS Position
// @Description Returns the ISS position for a specific timestamp (within 4 hours back/forward)
// @Tags ISS
// @Accept json
// @Produce json
// @Param timestamp path int true "Unix timestamp"
// @Param units query string false "Units (kilometers or miles)" Enums(kilometers, miles) default(kilometers)
// @Success 200 {object} models.ISSPosition
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /iss/historical/{timestamp} [get]
func (h *ISSHandler) GetHistoricalPosition(w http.ResponseWriter, r *http.Request) {
	timestampStr := chi.URLParam(r, "timestamp")
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid timestamp", "Timestamp must be a valid Unix timestamp")
		return
	}

	units := r.URL.Query().Get("units")
	if units != "miles" {
		units = "kilometers"
	}

	position, err := h.issService.GetHistoricalPosition(timestamp, units)
	if err != nil {
		if err.Error() == "timestamp outside retention window (4 hours back/forward)" {
			h.sendErrorResponse(w, http.StatusBadRequest, "Timestamp out of range", err.Error())
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get historical position", err.Error())
		return
	}

	if position == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Position not found", "No position data found for the specified timestamp")
		return
	}

	h.sendJSONResponse(w, http.StatusOK, position)
}

// GetPositionsInRange returns ISS positions within a time range
// @Summary Get ISS Positions in Time Range
// @Description Returns all ISS positions within a specified time range
// @Tags ISS
// @Accept json
// @Produce json
// @Param start_time query int true "Start timestamp (Unix)"
// @Param end_time query int true "End timestamp (Unix)"
// @Param units query string false "Units (kilometers or miles)" Enums(kilometers, miles) default(kilometers)
// @Success 200 {array} models.ISSPosition
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /iss/range [get]
func (h *ISSHandler) GetPositionsInRange(w http.ResponseWriter, r *http.Request) {
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Missing parameters", "Both start_time and end_time are required")
		return
	}

	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid start_time", "start_time must be a valid Unix timestamp")
		return
	}

	endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid end_time", "end_time must be a valid Unix timestamp")
		return
	}

	if startTime >= endTime {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid time range", "start_time must be less than end_time")
		return
	}

	if endTime-startTime > 24*3600 {
		h.sendErrorResponse(w, http.StatusBadRequest, "Time range too large", "Maximum time range is 24 hours")
		return
	}

	units := r.URL.Query().Get("units")
	if units != "miles" {
		units = "kilometers"
	}

	positions, err := h.issService.GetPositionsInRange(startTime, endTime, units)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get positions", err.Error())
		return
	}

	h.sendJSONResponse(w, http.StatusOK, positions)
}

// GetISSStatus returns general ISS tracking status and statistics
// @Summary Get ISS Tracking Status
// @Description Returns statistics about ISS tracking data and system status
// @Tags ISS
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /iss/status [get]
func (h *ISSHandler) GetISSStatus(w http.ResponseWriter, r *http.Request) {
	currentPos, err := h.issService.GetCurrentPosition("kilometers")
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get current position", err.Error())
		return
	}

	stats, err := h.issService.GetStatistics()
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get statistics", err.Error())
		return
	}

	status := map[string]interface{}{
		"status":              "operational",
		"current_position":    currentPos,
		"data_retention":      "4 hours back and forward",
		"collection_interval": "10 seconds",
		"last_update":         time.Unix(currentPos.Timestamp, 0).Format(time.RFC3339),
		"api_source":          "wheretheiss.at",
		"supported_units":     []string{"kilometers", "miles"},
		"statistics":          stats,
	}

	h.sendJSONResponse(w, http.StatusOK, status)
}

// PostHistoricalRequest handles POST request for historical data with JSON body
// @Summary Get Historical ISS Position (POST)
// @Description Returns the ISS position for a timestamp provided in request body
// @Tags ISS
// @Accept json
// @Produce json
// @Param request body models.HistoricalRequest true "Historical position request"
// @Success 200 {object} models.ISSPosition
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /iss/historical [post]
func (h *ISSHandler) PostHistoricalRequest(w http.ResponseWriter, r *http.Request) {
	var req models.HistoricalRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if req.Timestamp == 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "Missing timestamp", "Timestamp is required")
		return
	}

	units := req.Units
	if units != "miles" {
		units = "kilometers"
	}

	position, err := h.issService.GetHistoricalPosition(req.Timestamp, units)
	if err != nil {
		if err.Error() == "timestamp outside retention window (4 hours back/forward)" {
			h.sendErrorResponse(w, http.StatusBadRequest, "Timestamp out of range", err.Error())
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get historical position", err.Error())
		return
	}

	if position == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Position not found", "No position data found for the specified timestamp")
		return
	}

	h.sendJSONResponse(w, http.StatusOK, position)
}

func (h *ISSHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ISSHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := models.ErrorResponse{
		Error:   error,
		Message: message,
	}

	json.NewEncoder(w).Encode(errorResp)
}
