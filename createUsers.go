package main

import (
	"encoding/json"
	"errors"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"net/http"
	"strings"
	"time"
)

type AcceptEmail struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type UserShown struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

func (cfg *apiConfig) CreateUsers(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := AcceptEmail{}

	err := decoder.Decode(&params)
	if err != nil {
		Send(w, 500, nil, "Error decoding parameters")
		return
	}

	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		Send(w, 500, nil, "Error Hashing Password")
		return
	}

	params.Role, err = validateRole(params.Role)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	userDb, err := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.Password,
		Role:           params.Role,
	})
	if err != nil {
		Send(w, 404, nil, "Cannot Create User")
		return
	}

	respBody := UserShown{
		CreatedAt: userDb.CreatedAt,
		Name:      userDb.Name,
		Email:     userDb.Email,
		Role:      userDb.Role,
	}

	Send(w, 201, respBody, "User registered")

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
