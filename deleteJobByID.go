package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)

func (cfg *apiConfig) deleteJobByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)         
    idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusNotFound)
		return
	}


	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(403)
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(403)
		return
	}

	_, err = cfg.DB.DeleteJobByID(req.Context(), database.DeleteJobByIDParams{
		ID: id,
		EmployerID: ValidatedID,
	})

	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)


}



