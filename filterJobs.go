package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"database/sql"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)



func (cfg *apiConfig) FilterJobs(w http.ResponseWriter, req *http.Request) {
	location := req.URL.Query().Get("location")
	jobType := req.URL.Query().Get("type")
	salaryMinStr := req.URL.Query().Get("salary_min")
	salaryMaxStr := req.URL.Query().Get("salary_max")

	var err error
	salaryMin := float32(1)
	salaryMax := float32(1)
	

	if salaryMinStr != "" {
		conv, err := strconv.ParseFloat(salaryMinStr, 32)
		if err != nil {
			// handle error, e.g., bad input
			http.Error(w, "Invalid salary_min", http.StatusBadRequest)
			return
		}
		salaryMin = float32(conv)
	}

	if salaryMaxStr != "" {
		conv, err := strconv.ParseFloat(salaryMaxStr, 32)
		if err != nil {
			http.Error(w, "Invalid salary_max", http.StatusBadRequest)
			return
		}
		salaryMax = float32(conv)
	}	


	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	_, err = auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(401)
		return
	}


	JobDb, err := cfg.DB.FiltersJobs(req.Context(), database.FiltersJobsParams{
		Column1:  sql.NullString{
			String: location,
			Valid:  true,
		},
		Column2: jobType,
		Column3: salaryMin,
		Column4: salaryMax,
	})
	if err != nil {
    	http.Error(w, "Cannot Retrieve Filtered Jobs", http.StatusNotFound)
        return
	}

	respBody := make([]jobs, 0, len(JobDb))
    for _, dbjob := range JobDb {
        jobResp := jobs{
           	ID: dbjob.ID,
			CreatedAt: dbjob.CreatedAt,
			UpdatedAt: dbjob.UpdatedAt,
			Title: dbjob.Title,
			Description: dbjob.Description,
			Location: dbjob.Location,
			Type: dbjob.Type,
			Salary: dbjob.Salary,
			EmployerID: dbjob.EmployerID,
        }
        respBody = append(respBody, jobResp)
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