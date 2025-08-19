package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

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

	Body, err := cfg.DB.GetApplyJobs(req.Context(), id)
	if err != nil {
		Send(w, 500, nil, "Cannot get Application")
		return
	}

	Jobs, err := cfg.DB.GetJobsByID(req.Context(), Body.JobID)
	if err != nil {
		Send(w, 404, nil, "Cannot Retrieve Job")
		return
	}

	if (ValidatedID != Jobs.EmployerID) && (ValidatedID != Body.ApplicantID) {
		Send(w, 404, nil, "Only Job Creators and Applicants Can See application of a Job")
		return
	}

	Send(w, 200, Body, "Application retrieved")

}
