package main

import (
	"fmt"
	"github.com/xaaaaaanny/Chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cant get token from header")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Token does not exist or expired: %v", err))
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.secret, 60*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant make jwt token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	})
}
