package handlers

import (
	"encoding/json"
	"net/http"

	"iss-model-backend/internal/services"
	"iss-model-backend/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Admin Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Admin credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /admin/login [post]
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, LoginResponse{Token: token})
}

// HandleRegister (opcjonalne, do stworzenia pierwszego admina)
// @Summary Register Admin
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body RegisterRequest true "Admin credentials"
// @Success 201 {object} map[string]string
// @Router /admin/register [post]
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	_, err := h.authService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}
