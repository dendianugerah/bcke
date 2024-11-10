package auth

import (
	"encoding/json"
	"net/http"

	"github.com/dendianugerah/bcke/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=LoginResponse} "Successfully logged in"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Invalid credentials"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusBadRequest, "invalid request", nil))
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	json.NewEncoder(w).Encode(response.NewResponse(http.StatusOK, "login successful", resp))
} 