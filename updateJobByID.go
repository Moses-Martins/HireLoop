package main

import (
	"encoding/json"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

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
		Send(w, 404, nil, "Cannot Retrieve Job")
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
		Send(w, 400, nil, "Only Employers can Update a Job")
		return
	}

	if respBodyInitial.ID != getJobs.EmployerID {
		Send(w, 400, nil, "You cannot update a job if you didn't create")
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
		Send(w, 404, nil, "Cannot Update Job")
		return
	}

	respBody, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		Send(w, 404, nil, "Cannot Retrieve Job")
		return
	}

	Send(w, 200, respBody, "Job updated")

}
