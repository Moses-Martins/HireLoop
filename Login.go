package main

import (
	"encoding/json"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AcceptsEmail struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

type UserDisplayed struct {
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Token        string    `json:"authorization_token"`
	RefreshToken string    `json:"refresh_token"`
}

// Login godoc
// @Summary Login with email and password
// @Description Authenticates a user using email and password and returns JWT + refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body AcceptsEmail true "Email and password"
// @Success 200 {object} UserDisplayed "Login successful, returns JWT and refresh token"
// @Failure 401 {object} map[string]string "Incorrect password"
// @Failure 404 {object} map[string]string "Email not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/login [post]
func (cfg *apiConfig) Login(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := AcceptsEmail{}

	err := decoder.Decode(&params)
	if err != nil {
		Send(w, 500, nil, "Error decoding parameters")
		return
	}

	respBodyInitial, err := cfg.DB.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		Send(w, 404, nil, "Incorrect email (Email cannot be found)")
		return
	}

	err = auth.CheckPasswordHash(params.Password, respBodyInitial.HashedPassword)
	if err != nil {
		Send(w, 401, nil, "Incorrect password")
		return
	}

	token, err := auth.MakeJWT(respBodyInitial.ID, cfg.JwtSecret, time.Duration(43200)*time.Second)
	if err != nil {
		Send(w, 500, nil, "Cannot generate token")
		return
	}

	refreshtoken, err := auth.MakeRefreshToken()
	if err != nil {
		Send(w, 500, nil, "Error generating Refresh Token")
		return
	}

	_, err = cfg.DB.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshtoken,
		UserID:    respBodyInitial.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	response := UserDisplayed{
		Name:         respBodyInitial.Name,
		CreatedAt:    respBodyInitial.CreatedAt,
		UpdatedAt:    respBodyInitial.UpdatedAt,
		Email:        respBodyInitial.Email,
		Role:         respBodyInitial.Role,
		Token:        token,
		RefreshToken: refreshtoken,
	}

	Send(w, 200, response, "Login Successful")
}
