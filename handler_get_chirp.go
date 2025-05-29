package main

import (
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDstr := r.PathValue("chirpID")

	if chirpIDstr == "" {
		respondWithError(w, http.StatusNotFound, "ID is empty")
		return
	}

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can`t parse string to uuid")
		return
	}

	chirpDB, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can`t get chirp")
		return
	}

	chirp := Chirp{
		ID:        chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body:      chirpDB.Body,
		UserID:    chirpDB.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
