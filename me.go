package main

import (
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type meStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

// Me godoc
// @Summary Get current user info
// @Description Returns details of the currently authenticated user
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} meStruct "User details retrieved"
// @Failure 401 {object} map[string]string "Invalid or missing token"
// @Router /auth/me [get]
func (cfg *apiConfig) Me(w http.ResponseWriter, req *http.Request) {

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

	userDb, err := cfg.DB.GetUserByID(req.Context(), ValidatedID)

	respBody := meStruct{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Name:      userDb.Name,
		Email:     userDb.Email,
		Role:      userDb.Role,
	}

	Send(w, 200, respBody, "User details retrieved")
}