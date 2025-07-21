package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/Moses-Martins/HireLoop/internal/auth"
)



func (cfg *apiConfig) getAllApp(w http.ResponseWriter, req *http.Request) {
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

	Body, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		http.Error(w, "Cannot Retrieve Job", http.StatusNotFound)
		return
	}

	if ValidatedID != Body.EmployerID {
		http.Error(w, "Only Job Creators Can See Applicants of a Job", http.StatusNotFound)
		return
	}

	Applicants, err := cfg.DB.GetApplyByJobID(req.Context(), Body.ID)
	if err != nil {
		http.Error(w, "Cannot Retrieve All Applicants", http.StatusNotFound)
		return
	}


	respBody := make([]applyJob, 0, len(Applicants))
    for _, app := range Applicants {
        apps := applyJob{
			ID: app.ID,
			ApplicantID: app.ApplicantID,
			JobID: app.JobID,
			ResumeUrl: app.ResumeUrl,
			Status: app.Status,
        }
        respBody = append(respBody, apps)
    }


	data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

}