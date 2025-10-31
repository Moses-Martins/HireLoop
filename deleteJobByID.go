package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)


// deleteJobByID godoc
// @Summary Delete a job
// @Description Delete a job by ID. Only the employer who created the job can delete it. Requires authentication.
// @Tags jobs
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Job ID (UUID)"
// @Success 204 {object} map[string]interface{} "Job deleted successfully"
// @Failure 400 {object} map[string]string "Invalid UUID"
// @Failure 401 {object} map[string]string "Invalid or missing token"
// @Failure 403 {object} map[string]string "Forbidden: only creator can delete job"
// @Failure 404 {object} map[string]string "Job not found"
// @Router /api/jobs/{id} [delete]
func (cfg *apiConfig) deleteJobByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		Send(w, 400, nil, "Invalid UUID")
		return
	}

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

	_, err = cfg.DB.DeleteJobByID(req.Context(), database.DeleteJobByIDParams{
		ID:         id,
		EmployerID: ValidatedID,
	})

	if err != nil {
		Send(w, 403, nil, "Forbidden")
		return
	}

	Send(w, 204, map[string]interface{}{}, "Job deleted")
}
