package main

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"time"

	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)

var userRole string

// Handler to start OAuth2 Register
func (cfg *apiConfig) GoogleRegister(w http.ResponseWriter, req *http.Request) {
	cfg.GoogleOauthConfig.RedirectURL = cfg.RegisterRedirectUrl
	userRole = req.URL.Query().Get("role")
	url := cfg.GoogleOauthConfig.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

// Handler to handle the callback from Google
func (cfg *apiConfig) RegisterCallback(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	code := req.URL.Query().Get("code")

	token, err := cfg.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		Send(w, 500, nil, "Failed to exchange token: "+err.Error())
		return
	}

	client := cfg.GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		Send(w, 500, nil, "Failed to get user info: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		Send(w, 500, nil, "Failed to decode user info: "+err.Error())
		return
	}

	hashedPass, err := auth.HashPassword("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	if err != nil {
		Send(w, 500, nil, "Error Hashing Password")
		return
	}

	if userRole == "" {
		Send(w, 400, nil, "Please specify the role you want to register for.")
		return
	}

	userRole, err = validateRole(userRole)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	userDb, err := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{
		Name:           userInfo.Name,
		Email:          userInfo.Email,
		HashedPassword: hashedPass,
		Role:           userRole,
	})
	if err != nil {
		Send(w, 404, nil, "Cannot Create User")
		return
	}

	loginToken, err := auth.MakeJWT(userDb.ID, cfg.JwtSecret, time.Duration(3600)*time.Second)
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
		UserID:    userDb.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	respBody := UserDisplayed{
		Name:         userDb.Name,
		CreatedAt:    userDb.CreatedAt,
		UpdatedAt:    userDb.UpdatedAt,
		Email:        userDb.Email,
		Role:         userDb.Role,
		Token:        loginToken,
		RefreshToken: refreshtoken,
	}

	Send(w, 200, respBody, "Google registration completed")

}
