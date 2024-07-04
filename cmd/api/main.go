package main

import (
	"github.com.mloughton/crud/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	serverHTTP, err := server.NewServer()
	if err != nil {
		panic(err)
	}

	serverHTTP.ListenAndServe()

}
