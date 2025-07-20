package main


import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/google/uuid" 
	"github.com/Moses-Martins/HireLoop/internal/auth"
)


func (cfg *apiConfig) getJobByID(w http.ResponseWriter, req *http.Request) {

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


	vars := mux.Vars(req)         
    idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusNotFound)
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