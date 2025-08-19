package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

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

	Body, err := cfg.DB.GetApplyJobs(req.Context(), id)
	if err != nil {
		Send(w, 500, nil, "Cannot get Application")
		return
	}

	if ValidatedID != Body.ApplicantID {
		Send(w, 404, nil, "Only the applicant can withdraw their job application")
		return
	}

	_, err = cfg.DB.DeleteAppByID(req.Context(), database.DeleteAppByIDParams{
		ID:          Body.ID,
		ApplicantID: Body.ApplicantID,
	})

	if err != nil {
		Send(w, 400, nil, "The request cannot be processed")
		return
	}

	Send(w, 204, map[string]interface{}{}, "Application deleted")

}
