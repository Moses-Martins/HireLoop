package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"net/http"
)

func (cfg *apiConfig) getAllJobs(w http.ResponseWriter, req *http.Request) {
	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	_, err = auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	JobDb, err := cfg.DB.GetAllJobs(req.Context())
	if err != nil {
		Send(w, 404, nil, "Cannot Retrieve Jobs")
		return
	}

	respBody := make([]jobs, 0, len(JobDb))
	for _, dbjob := range JobDb {
		jobResp := jobs{
			ID:          dbjob.ID,
			CreatedAt:   dbjob.CreatedAt,
			UpdatedAt:   dbjob.UpdatedAt,
			Title:       dbjob.Title,
			Description: dbjob.Description,
			Location:    dbjob.Location,
			Type:        dbjob.Type,
			Salary:      dbjob.Salary,
			EmployerID:  dbjob.EmployerID,
		}
		respBody = append(respBody, jobResp)
	}

	Send(w, 200, respBody, "Jobs retrieved")

}
