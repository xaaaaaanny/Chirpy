package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type returnVals struct {
	Body        string `json:"body"`
	Error       string `json:"error"`
	CleanedBody string `json:"cleaned_body"`
}

func handlerChirp(w http.ResponseWriter, r *http.Request) {

	respBody := returnVals{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&respBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}

	if len(respBody.Body) > 140 {
		respBody.Error = "Chirp is too long"

		respondWithError(w, http.StatusBadRequest, respBody.Error)

		return
	}

	respBody.CleanedBody = clearText(respBody.Body)
	respondWithJSON(w, http.StatusOK, respBody)
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

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	resp := returnVals{
		Error: message,
	}

	respondWithJSON(w, statusCode, resp)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
