package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/Moses-Martins/HireLoop/internal/auth"
)



func (cfg *apiConfig) viewSingleApp(w http.ResponseWriter, req *http.Request) {
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

	Jobs, err := cfg.DB.GetJobsByID(req.Context(), Body.JobID)
	if err != nil {
		http.Error(w, "Cannot Retrieve Job", http.StatusNotFound)
		return
	}

	if (ValidatedID != Jobs.EmployerID) && (ValidatedID != Body.ApplicantID) {
		http.Error(w, "Only Job Creators and Applicants Can See application of a Job", http.StatusNotFound)
		return
	}
 
	data, err := json.Marshal(Body)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)


}