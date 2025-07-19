package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/Moses-Martins/HireLoop/internal/auth"
)


type meStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name string `json:"name"`
	Email     string    `json:"email"`
	Role string `json:"role"`
}


func (cfg *apiConfig) Me(w http.ResponseWriter, req *http.Request) {
	
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

	userDb, err := cfg.DB.GetUserByID(req.Context(), ValidatedID)


	respBody := meStruct{
		ID:       userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Name:	userDb.Name,
		Email:    userDb.Email,
		Role: userDb.Role,
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