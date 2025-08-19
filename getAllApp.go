package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

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

	Body, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		Send(w, 404, nil, "Cannot Retrieve Job")
		return
	}

	if ValidatedID != Body.EmployerID {
		Send(w, 404, nil, "Only Job Creators Can See Applicants of a Job")
		return
	}

	Applicants, err := cfg.DB.GetApplyByJobID(req.Context(), Body.ID)
	if err != nil {
		Send(w, 404, nil, "Cannot Retrieve All Applicants")
		return
	}

	respBody := make([]applyJob, 0, len(Applicants))
	for _, app := range Applicants {
		apps := applyJob{
			ID:          app.ID,
			ApplicantID: app.ApplicantID,
			JobID:       app.JobID,
			ResumeUrl:   app.ResumeUrl,
			Status:      app.Status,
		}
		respBody = append(respBody, apps)
	}

	Send(w, 200, respBody, "Applications retrieved")

}
