package server

import (
	"log"
	"net/http"

	"github.com.mloughton/crud/internal/dax"
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
	columns, err := dax.ParseInput(r.Form.Get("input"))
	if err != nil {
		// respondWithError(w, http.StatusBadRequest, "Bad Request")
		w.Header().Add("HX-Retarget", "#error")
		w.Write([]byte(err.Error()))
		log.Print(err)
		return
	}
	dax, err := dax.GenerateDax(columns)
	if err != nil {
		log.Print(err)
		return
	}
	w.Write([]byte(dax))
}
