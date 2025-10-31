package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// getAllApp godoc
// @Summary Get all applications for a job
// @Description Retrieve all applicants for a specific job. Only the employer who created the job can access this. Requires authentication.
// @Tags applications
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Job ID (UUID)"
// @Success 200 {array} applyJob "List of applications"
// @Failure 401 {object} map[string]string "Missing or invalid token"
// @Failure 404 {object} map[string]string "Job not found or access forbidden"
// @Router /api/employers/{id}/applications [get]
func (cfg *apiConfig) getAllApp(w http.ResponseWriter, req *http.Request) {
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

	jobBody, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		Send(w, 404, nil, "Cannot retrieve job")
		return
	}

	if ValidatedID != jobBody.EmployerID {
		Send(w, 404, nil, "Only job creators can see applicants for this job")
		return
	}

	applicants, err := cfg.DB.GetApplyByJobID(req.Context(), jobBody.ID)
	if err != nil {
		Send(w, 404, nil, "Cannot retrieve applicants")
		return
	}

	respBody := make([]applyJob, 0, len(applicants))
	for _, app := range applicants {
		respBody = append(respBody, applyJob{
			ID:          app.ID,
			ApplicantID: app.ApplicantID,
			JobID:       app.JobID,
			ResumeUrl:   app.ResumeUrl,
			Status:      app.Status,
		})
	}

	Send(w, 200, respBody, "Applications retrieved")
}

