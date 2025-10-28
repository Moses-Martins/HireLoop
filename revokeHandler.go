package main

import(
	"net/http"
	"time"
	"github.com/Moses-Martins/HireLoop/internal/auth"
)


func (cfg *apiConfig) RevokeHandler(w http.ResponseWriter, req *http.Request) {
	
	token_string, err := auth.GetBearerToken(req.Header)
    if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("missing authorization header"))
		return
	}

	RefreshTokenDb, err := cfg.DB.GetRefreshTokens(req.Context())
	if err != nil {
    	http.Error(w, "Cannot Retrieve Refresh Tokens", http.StatusNotFound)
        return
	}


	dbToStruct := make([]RefreshStruct, 0, len(RefreshTokenDb))
    for _, dbToken := range RefreshTokenDb {
        RefreshResp := RefreshStruct{
            Token:        dbToken.Token,
            CreatedAt: dbToken.CreatedAt,
            UpdatedAt: dbToken.UpdatedAt,
            UserID:      dbToken.UserID,
			ExpiresAt: dbToken.ExpiresAt,
        }
        dbToStruct = append(dbToStruct, RefreshResp)

    }

	respBodyInitial, exist := findRefreshTokenByToken(dbToStruct, token_string)
	if exist != true {
		w.WriteHeader(401)
		return
	}

	if respBodyInitial.ExpiresAt.Before(time.Now()) {
		w.WriteHeader(401)
		return
	}

    cfg.DB.RevokeToken(req.Context(), respBodyInitial.Token) 

	w.WriteHeader(204)

}