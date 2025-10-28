package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// deleteApplication godoc
// @Summary Withdraw a job application
// @Description Allows an applicant to delete their own job application. Requires authentication.
// @Tags applications
// @Security ApiKeyAuth
// @Param id path string true "Application ID (UUID)"
// @Success 204 {object} map[string]interface{} "Application deleted"
// @Failure 401 {object} map[string]string "Missing or invalid token"
// @Failure 404 {object} map[string]string "Application not found or unauthorized"
// @Failure 400 {object} map[string]string "Cannot process request"
// @Router /applications/{id} [delete]
func (cfg *apiConfig) deleteApplication(w http.ResponseWriter, req *http.Request) {
	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
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

	application, err := cfg.DB.GetApplyJobs(req.Context(), id)
	if err != nil {
		Send(w, 500, nil, "Cannot get application")
		return
	}

	if ValidatedID != application.ApplicantID {
		Send(w, 404, nil, "Only the applicant can withdraw their job application")
		return
	}

	_, err = cfg.DB.DeleteAppByID(req.Context(), database.DeleteAppByIDParams{
		ID:          application.ID,
		ApplicantID: application.ApplicantID,
	})

	if err != nil {
		Send(w, 400, nil, "The request cannot be processed")
		return
	}

	Send(w, 204, map[string]interface{}{}, "Application deleted")
}
