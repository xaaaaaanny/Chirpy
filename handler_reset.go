package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You have no permission"))
		return
	}
	cfg.db.ResetUsers(context.Background())
	cfg.db.ResetChirps(context.Background())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0.\nUsers and chirps deleted\n"))
}
