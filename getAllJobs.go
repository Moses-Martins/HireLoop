package main

import (
	"net/http"
	"encoding/json"
	"log"
)



func (cfg *apiConfig) getAllJobs(w http.ResponseWriter, req *http.Request) {

	JobDb, err := cfg.DB.GetAllJobs(req.Context())
	if err != nil {
    	http.Error(w, "Cannot Retrieve Chirps", http.StatusNotFound)
        return
	}

	respBody := make([]jobs, 0, len(JobDb))
    for _, dbjob := range JobDb {
        jobResp := jobs{
           	ID: dbjob.ID,
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