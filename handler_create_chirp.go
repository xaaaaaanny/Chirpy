package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/xaaaaaanny/Chirpy/internal/auth"
	"github.com/xaaaaaanny/Chirpy/internal/database"
	"log"
	"net/http"
	"strings"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("invalid token in header: %v", err))
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("JWT is invalid: %v", err))
		return
	}

	type Params struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id)"`
	}
	params := &Params{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	params.Body = clearText(params.Body)

	arg := database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	}

	chirpDB, err := cfg.db.CreateChirp(r.Context(), arg)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Can`t create chirp")
		return
	}

	chirpJSON := Chirp{
		ID:        chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body:      chirpDB.Body,
		UserID:    chirpDB.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirpJSON)
}

func clearText(text string) string {
	splitText := strings.Fields(text)

	for i, word := range splitText {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			splitText[i] = "****"
		}
	}
	return strings.Join(splitText, " ")
}
