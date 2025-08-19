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

// Handler to start OAuth2 login
func (cfg *apiConfig) GoogleLogin(w http.ResponseWriter, req *http.Request) {
	cfg.GoogleOauthConfig.RedirectURL = cfg.LoginRedirectUrl
	url := cfg.GoogleOauthConfig.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

// Handler to handle the callback from Google
func (cfg *apiConfig) LoginCallback(w http.ResponseWriter, req *http.Request) {
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

	Information, err := cfg.DB.GetUserByEmail(req.Context(), userInfo.Email)
	if err != nil {
		Send(w, 404, nil, "User not found. Please register to continue.")
		return
	}

	loginToken, err := auth.MakeJWT(Information.ID, cfg.JwtSecret, time.Duration(43200)*time.Second)
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
		UserID:    Information.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	respBody := UserDisplayed{
		Name:         Information.Name,
		CreatedAt:    Information.CreatedAt,
		UpdatedAt:    Information.UpdatedAt,
		Email:        Information.Email,
		Role:         Information.Role,
		Token:        loginToken,
		RefreshToken: refreshtoken,
	}

	Send(w, 200, respBody, "Google login completed")

}
