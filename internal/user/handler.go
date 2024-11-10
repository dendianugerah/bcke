package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/dendianugerah/bcke/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create new user
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User registration details"
// @Success 201 {object} response.Response{data=User} "User created successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /register [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusBadRequest, "invalid request", nil))
		return
	}

	user, err := h.service.Create(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.NewResponse(http.StatusCreated, "user created", user))
}

// List godoc
// @Summary List users
// @Description Get list of users with pagination, sorting and filtering
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Items per page"
// @Param sort query string false "Sort field"
// @Param search query string false "Search term"
// @Security Bearer
// @Success 200 {object} response.Response{data=[]User} "Users retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /users [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = 10
	}

	filter := FilterOptions{
		Page:     page,
		PageSize: pageSize,
		Sort:     r.URL.Query().Get("sort"),
		Search:   r.URL.Query().Get("search"),
	}

	users, err := h.service.List(r.Context(), filter)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	json.NewEncoder(w).Encode(response.NewResponse(http.StatusOK, "users retrieved", users))
}

// Update godoc
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User update details"
// @Security Bearer
// @Success 200 {object} response.Response{data=User} "User updated successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /users/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusBadRequest, "invalid request", nil))
		return
	}

	user, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	json.NewEncoder(w).Encode(response.NewResponse(http.StatusOK, "user updated", user))
}

// Delete godoc
// @Summary Delete user
// @Description Delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security Bearer
// @Success 200 {object} response.Response "User deleted successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /users/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.Delete(r.Context(), id); err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	json.NewEncoder(w).Encode(response.NewResponse(http.StatusOK, "user deleted", nil))
} 