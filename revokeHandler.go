package main

import(
	"net/http"
	"time"
	"github.com/Moses-Martins/HireLoop/internal/auth"
)


// RevokeHandler godoc
// @Summary Revoke a refresh token
// @Description Revokes the provided refresh token from the Authorization header, effectively logging the user out
// @Tags auth
// @Security ApiKeyAuth
// @Success 204 "Refresh token revoked successfully"
// @Failure 401 {object} map[string]string "Invalid, missing, or expired token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /revoke [post]
func (cfg *apiConfig) RevokeHandler(w http.ResponseWriter, req *http.Request) {

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

	if err := cfg.DB.RevokeToken(req.Context(), respBodyInitial.Token); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
