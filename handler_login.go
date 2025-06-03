package main

import (
	"encoding/json"
	"github.com/xaaaaaanny/Chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email    string
		Password string
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

	cfg.user_id = userDB.ID

	respondWithJSON(w, http.StatusOK, User{
		ID:        userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email:     userDB.Email,
	})
}
