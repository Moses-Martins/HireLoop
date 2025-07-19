package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"errors"
	"strings"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)

type AcceptEmail struct {
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
	Role string `json:"role"`
}
 
type UserShown struct {
	CreatedAt time.Time `json:"created_at"`
	Name string `json:"name"`
	Email     string    `json:"email"`
	Role string `json:"role"`
}

func (cfg *apiConfig) CreateUsers(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := AcceptEmail{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error Hashing Password: %s", err)
		w.WriteHeader(500)
		return
	}

	params.Role, err = validateRole(params.Role)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	userDb, err := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{		
		Name:  params.Name,  
		Email:	params.Email,    
		HashedPassword: params.Password,
		Role:	params.Role,
	})
	if err != nil {
    	http.Error(w, "Cannot Create User", http.StatusNotFound)
        return
	}

	respBody := UserShown{
		CreatedAt: userDb.CreatedAt,
		Name: userDb.Name,
		Email: userDb.Email,
		Role: userDb.Role,
	}

	data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
	
}






func validateRole(role string) (string, error) {
    normalized := strings.ToLower(strings.TrimSpace(role))

    switch normalized {
    case "employer", "applicant":
        return normalized, nil
    default:
        return "", errors.New("invalid role: must be 'employer' or 'applicant'")
    }
}