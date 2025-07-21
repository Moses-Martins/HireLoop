package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)



func (cfg *apiConfig) deleteApplication(w http.ResponseWriter, req *http.Request) {
	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	vars := mux.Vars(req)         
    idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusNotFound)
		return
	}

	Body, err := cfg.DB.GetApplyJobs(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot get Application"))
		return
	}

	if ValidatedID != Body.ApplicantID {
		http.Error(w, "Only the applicant can withdraw their job application", http.StatusNotFound)
		return
	}

	_, err = cfg.DB.DeleteAppByID(req.Context(), database.DeleteAppByIDParams{
		ID: Body.ID,
		ApplicantID: Body.ApplicantID,
	})

	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
 
}