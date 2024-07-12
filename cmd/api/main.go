package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mloughton/pbifiltergen/internal/server"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	serverHTTP, err := server.NewServer()
	if err != nil {
		panic(err)
	}

	log.Fatal(serverHTTP.ListenAndServe())

}
