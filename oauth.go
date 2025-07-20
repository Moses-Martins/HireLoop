package main

import (
	"context"
	"encoding/json"
	"net/http"
	"log"

	"golang.org/x/oauth2"
	"time"

	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)


// Handler to start OAuth2 login
func (cfg *apiConfig) GoogleLogin(w http.ResponseWriter, req *http.Request) {
	url := cfg.GoogleOauthConfig.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

// Handler to handle the callback from Google
func (cfg *apiConfig) GoogleCallback(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	code := req.URL.Query().Get("code")
	token, err := cfg.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := cfg.GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPass, err := auth.HashPassword("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	if err != nil {
		log.Printf("Error Hashing Password: %s", err)
		w.WriteHeader(500)
		return
	}
	

	Information, err := cfg.DB.GetUserByEmail(req.Context(), userInfo.Email)
	if err != nil {
    	userDb, err := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{		
			Name:  userInfo.Name,  
			Email:	userInfo.Email,
			HashedPassword: hashedPass,    
			Role:	"applicant",
		})
		if err != nil {
			http.Error(w, "Cannot Create User", http.StatusNotFound)
			return
		}

		loginToken, err := auth.MakeJWT(userDb.ID, cfg.JwtSecret, time.Duration(3600) * time.Second)
		if err != nil {
			log.Printf("Cannot generate token %s", err)
			w.WriteHeader(500)
			return
		}

		refreshtoken, err := auth.MakeRefreshToken()
		if err != nil {
			log.Printf("Error generating Refresh Token: %s", err)
			w.WriteHeader(500)
			return
		}

		_, err = cfg.DB.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
			Token: refreshtoken,
			UserID: userDb.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		})

		respBody := UserDisplayed{
			Name: userDb.Name,
			Email: userDb.Email,
			Role: userDb.Role,
			Token: loginToken,
			RefreshToken: refreshtoken,
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
		return
	}

	loginToken, err := auth.MakeJWT(Information.ID, cfg.JwtSecret, time.Duration(3600) * time.Second)
	if err != nil {
		log.Printf("Cannot generate token %s", err)
		w.WriteHeader(500)
		return
	}

	refreshtoken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error generating Refresh Token: %s", err)
		w.WriteHeader(500)
		return
	}

	_, err = cfg.DB.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshtoken,
		UserID: Information.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	respBody := UserDisplayed{
		Name: Information.Name,
		Email: Information.Email,
		Role: Information.Role,
		Token: loginToken,
		RefreshToken: refreshtoken,
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