package main

import (
	"encoding/json"
	"fmt"
	"github.com/xaaaaaanny/Chirpy/internal/auth"
	"github.com/xaaaaaanny/Chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := &Params{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	userDB, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cant get user by email")
		return
	}

	if auth.CheckPasswordHash(userDB.HashedPassword, params.Password) != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(userDB.ID, cfg.secret, 60*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("can`t make JWT: %v", err))
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    userDB.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can`t create refresh token")
	}
	//cfg.user_id = userDB.ID

	respondWithJSON(w, http.StatusOK, User{
		ID:           userDB.ID,
		CreatedAt:    userDB.CreatedAt,
		UpdatedAt:    userDB.UpdatedAt,
		Email:        userDB.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}
