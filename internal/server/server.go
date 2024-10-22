package server

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"
)

type Server struct {
	limiter *rate.Limiter
}

func NewServer() (*http.Server, error) {
	host := ""
	local := flag.Bool("local", false, "Enable local mode")
	flag.Parse()

	if *local {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("PORT environment variable not set")
	}

	server := Server{
		limiter: rate.NewLimiter(rate.Every(500*time.Millisecond), 5),
	}

	serverHTTP := &http.Server{
		Addr:              host + ":" + port,
		Handler:           server.RegisterRoutes(),
		ReadHeaderTimeout: time.Minute,
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
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Fatal(err)
	}
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
