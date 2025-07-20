package main


import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/google/uuid" 
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)



func (cfg *apiConfig) updateJobByID(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	params := job{}
	
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	params.Salary, err = validateFloat(params.Salary)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	params.Type, err = validateJobType(params.Type)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	vars := mux.Vars(req)         
    idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusNotFound)
		return
	}

	getJobs, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		http.Error(w, "Cannot Retrieve Job", http.StatusNotFound)
		return
	}

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

	respBodyInitial, err := cfg.DB.GetUserByID(req.Context(), ValidatedID)
	if err != nil {
    	w.WriteHeader(401)
		return
	}

	if respBodyInitial.Role != "employer" {
		w.WriteHeader(400)
		w.Write([]byte("Only Employers can Update a Job"))
		return
	}


	if respBodyInitial.ID != getJobs.EmployerID {
		w.WriteHeader(400)
		w.Write([]byte("You cannot update a job if you didn't create"))
		return
	}


	err = cfg.DB.UpdateJob(req.Context(), database.UpdateJobParams{
		Title: params.Title,
		Description: params.Description,
		Location:   params.Location,
		Type:        params.Type,
		Salary:      params.Salary,
		ID:          id,
	})
	if err != nil {
    	http.Error(w, "Cannot Update Job", http.StatusNotFound)
        return
	}


	respBody, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		http.Error(w, "Cannot Retrieve Job", http.StatusNotFound)
		return
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