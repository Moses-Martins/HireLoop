package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// viewSingleApp godoc
// @Summary Get a single job application
// @Description Retrieve details of a single job application. Accessible only by the applicant or the employer who posted the job. Requires authentication.
// @Tags applications
// @Security ApiKeyAuth
// @Param id path string true "Application ID (UUID)"
// @Success 200 {object} applyJob "Application details"
// @Failure 401 {object} map[string]string "Missing or invalid token"
// @Failure 404 {object} map[string]string "Application not found or unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /applications/{id} [get]
func (cfg *apiConfig) viewSingleApp(w http.ResponseWriter, req *http.Request) {
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

	job, err := cfg.DB.GetJobsByID(req.Context(), application.JobID)
	if err != nil {
		Send(w, 404, nil, "Cannot retrieve job")
		return
	}

	if (ValidatedID != job.EmployerID) && (ValidatedID != application.ApplicantID) {
		Send(w, 404, nil, "Only job creators and applicants can view this application")
		return
	}

	Send(w, 200, application, "Application retrieved")
}
