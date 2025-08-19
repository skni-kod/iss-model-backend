package handlers

import (
	"net/http"

	"iss-model-backend/internal/services"
	"iss-model-backend/internal/utils"
)

type CrewHandler struct {
	crewService *services.CrewService
}

func NewCrewHandler(crewService *services.CrewService) *CrewHandler {
	return &CrewHandler{
		crewService: crewService,
	}
}

// GetCurrentCrew returns the current ISS crew
// @Summary Get current ISS crew
// @Description Return current crew aboard the ISS
// @Tags Crew
// @Accept json
// @Success 200 {object} models.ISSCrewResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /iss/crew [get]
func (h *CrewHandler) GetCurrentCrew(w http.ResponseWriter, r *http.Request) {
	crew, err := h.crewService.GetCurrentCrew()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get current crew", err.Error())
	}

	utils.SendJSONResponse(w, http.StatusOK, crew)
}

// GetCurrentCrewWithPhotos returns the current ISS crew with photos
// @Summary Get current ISS crew with photos
// @Description Return current crew aboard the ISS with photos from wikipedia API
// @Tags Crew
// @Accept json
// @Success 200 {object} models.ISSCrewWithPhotosResponse
// @Failure 500 {object} models.ErrorResponse
func (h *CrewHandler) GetCurrentCrewWithPhotos(w http.ResponseWriter, r *http.Request) {
	crewWithPhotos, err := h.crewService.GetCurrentCrewWithPhotos()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get current crew", err.Error())
	}

	utils.SendJSONResponse(w, http.StatusOK, crewWithPhotos)
}
