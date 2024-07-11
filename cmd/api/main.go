package main

import (
	"log"

	"github.com.mloughton/crud/internal/server"
	"github.com/joho/godotenv"
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
