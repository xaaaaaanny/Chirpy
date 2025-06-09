package main

import (
	"github.com/xaaaaaanny/Chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant get token")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant revoke refresh token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
