package main

import (
	"encoding/json"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// updateJobByID godoc
// @Summary Update a job
// @Description Update an existing job. Only the employer who created the job can update it. Requires authentication.
// @Tags jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Job ID (UUID)"
// @Param job body job true "Job details to update"
// @Success 200 {object} jobs "Job updated successfully"
// @Failure 400 {object} map[string]string "Invalid input or unauthorized"
// @Failure 401 {object} map[string]string "Missing or invalid authentication token"
// @Failure 403 {object} map[string]string "Forbidden: only creator can update job"
// @Failure 404 {object} map[string]string "Job not found or cannot update"
// @Router /jobs/{id} [put]
func (cfg *apiConfig) updateJobByID(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	params := job{}

	err := decoder.Decode(&params)
	if err != nil {
		Send(w, 500, nil, "Error decoding parameters")
		return
	}

	params.Salary, err = validateFloat(params.Salary)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	params.Type, err = validateJobType(params.Type)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		Send(w, 404, nil, "Invalid UUID")
		return
	}

	getJobs, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		Send(w, 404, nil, "Cannot retrieve job")
		return
	}

	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		Send(w, 401, nil, "Missing or invalid authentication token")
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		Send(w, 401, nil, "Missing or invalid authentication token")
		return
	}

	respBodyInitial, err := cfg.DB.GetUserByID(req.Context(), ValidatedID)
	if err != nil {
		Send(w, 403, nil, "Forbidden")
		return
	}

	if respBodyInitial.Role != "employer" {
		Send(w, 400, nil, "Only employers can update a job")
		return
	}

	if respBodyInitial.ID != getJobs.EmployerID {
		Send(w, 400, nil, "You cannot update a job you didn't create")
		return
	}

	err = cfg.DB.UpdateJob(req.Context(), database.UpdateJobParams{
		Title:       params.Title,
		Description: params.Description,
		Location:    params.Location,
		Type:        params.Type,
		Salary:      params.Salary,
		ID:          id,
	})
	if err != nil {
		Send(w, 404, nil, "Cannot update job")
		return
	}

	respBody, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		Send(w, 404, nil, "Cannot retrieve job")
		return
	}

	Send(w, 200, respBody, "Job updated")
}
