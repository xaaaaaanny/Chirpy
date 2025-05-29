package main

import "net/http"

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []Chirp

	chirpsDB, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can`t get chirs from DB")
		return
	}

	if len(chirpsDB) == 0 {
		respondWithError(w, http.StatusInternalServerError, "No chirp was created yet")
		return
	}

	for _, chirp := range chirpsDB {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
