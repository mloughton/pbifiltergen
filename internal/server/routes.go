package server

import (
	"fmt"
	"log"
	"net/http"
)

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", GetHealthHandler)

	mux.Handle("/", http.FileServer(http.Dir("./cmd/web/assets")))

	mux.HandleFunc("POST /input", PostInputHandler)
	return mux
}

func GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	respondWithJSON(w, http.StatusOK, res)
}

func PostInputHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		log.Print(err)
		return
	}
	w.Write([]byte(fmt.Sprintf("recieved post: %s", r.Form.Get("input"))))
}
