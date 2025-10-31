package main

import (
	"log"
	"net/http"
	"time"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/Moses-Martins/HireLoop/internal/auth"
)

type RefreshStruct struct {
    Token string
    CreatedAt time.Time
    UpdatedAt time.Time
    UserID uuid.UUID
    ExpiresAt time.Time
}

type RefreshToken struct {
    Token string `json:"token"`
}


// RefreshHandler godoc
// @Summary Refresh JWT using refresh token
// @Description Accepts a valid refresh token in the Authorization header and returns a new JWT
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} RefreshToken "Returns new JWT token"
// @Failure 401 {object} map[string]string "Invalid, missing, or expired refresh token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/refresh [post]
func (cfg *apiConfig) RefreshHandler(w http.ResponseWriter, req *http.Request) {

	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("missing authorization header"))
		return
	}

	RefreshTokenDb, err := cfg.DB.GetRefreshTokens(req.Context())
	if err != nil {
		http.Error(w, "Cannot retrieve refresh tokens", http.StatusNotFound)
		return
	}

	dbToStruct := make([]RefreshStruct, 0, len(RefreshTokenDb))
	for _, dbToken := range RefreshTokenDb {
		RefreshResp := RefreshStruct{
			Token:     dbToken.Token,
			CreatedAt: dbToken.CreatedAt,
			UpdatedAt: dbToken.UpdatedAt,
			UserID:    dbToken.UserID,
			ExpiresAt: dbToken.ExpiresAt,
		}
		dbToStruct = append(dbToStruct, RefreshResp)
	}

	respBodyInitial, exist := findRefreshTokenByToken(dbToStruct, token_string)
	if !exist {
		w.WriteHeader(401)
		return
	}

	if respBodyInitial.ExpiresAt.Before(time.Now()) {
		w.WriteHeader(401)
		return
	}

	UserData, err := cfg.DB.GetUserByRefreshToken(req.Context(), respBodyInitial.Token)
	if err != nil {
		log.Printf("Cannot generate new token %s", err)
		w.WriteHeader(500)
		return
	}

	token, err := auth.MakeJWT(UserData.ID, cfg.JwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		log.Printf("Cannot generate token %s", err)
		w.WriteHeader(500)
		return
	}

	respBody := RefreshToken{
		Token: token,
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


func findRefreshTokenByToken(tokens []RefreshStruct, token string) (RefreshStruct, bool) {
	for _, t := range tokens {
		if t.Token == token {
			return t, true // found
		}
	}
	return RefreshStruct{}, false // not found
}
