package main

import (
	"net/http"
	"encoding/json"
	"log"
	"database/sql"
	"github.com/Moses-Martins/HireLoop/internal/auth"
)


func (cfg *apiConfig) searchJobs(w http.ResponseWriter, req *http.Request) {
	str := req.URL.Query().Get("keyword")

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


	JobDb, err := cfg.DB.SearchJobs(req.Context(), sql.NullString{
		String: str,
		Valid:  true,
	})
	if err != nil {
		http.Error(w, "Cannot search the database", http.StatusNotFound)
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
	return
		
	
}