package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/xaaaaaanny/Chirpy/internal/auth"
	"github.com/xaaaaaanny/Chirpy/internal/database"
	"net/http"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cant hash given password")
		return
	}

	userDB, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can`t create user")
		return
	}
	userJSON := User{
		ID:        userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email:     userDB.Email,
	}

	//cfg.user_id = userJSON.ID

	respondWithJSON(w, http.StatusCreated, userJSON)
}
