package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func (cfg *apiConfig) deleteJobByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		Send(w, 400, nil, "Invalid UUID")
		return
	}

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

	_, err = cfg.DB.DeleteJobByID(req.Context(), database.DeleteJobByIDParams{
		ID:         id,
		EmployerID: ValidatedID,
	})

	if err != nil {
		Send(w, 403, nil, "Forbidden")
		return
	}

	Send(w, 204, map[string]interface{}{}, "Job deleted")

}
