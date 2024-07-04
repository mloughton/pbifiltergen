package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Server struct {
}

func NewServer() (*http.Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("can't get PORT in env")
	}

	server := &Server{}
	serverHTTP := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: server.RegisterRoutes(),
	}
	return serverHTTP, nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	responseJSON, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJSON)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type response struct {
		Error string `json:"error"`
	}
	responseBody := response{
		Error: msg,
	}
	respondWithJSON(w, code, responseBody)
}
