package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// getJobByID godoc
// @Summary Get job by ID
// @Description Retrieve a single job listing by its ID. Requires authentication.
// @Tags jobs
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Job ID (UUID)"
// @Success 200 {object} jobs "Job retrieved successfully"
// @Failure 401 {object} map[string]string "Invalid or missing token"
// @Failure 404 {object} map[string]string "Job not found or invalid UUID"
// @Router /api/jobs/{id} [get]
func (cfg *apiConfig) getJobByID(w http.ResponseWriter, req *http.Request) {

	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	_, err = auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		Send(w, 404, nil, "Invalid UUID")
		return
	}

	respBody, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		Send(w, 404, nil, "Cannot retrieve job")
		return
	}

	Send(w, 200, respBody, "Job retrieved")
}