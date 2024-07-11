package server

import (
	"log"
	"net/http"

	"github.com/mloughton/pbifiltergen/cmd/web/staticfs"
	"github.com/mloughton/pbifiltergen/internal/dax"
)

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", GetHealthHandler)

	mux.HandleFunc("/", AppHandler)

	mux.HandleFunc("POST /input", PostInputHandler)

	mux.HandleFunc("GET /copy", GetCopyHandler)
	return mux
}

func AppHandler(w http.ResponseWriter, r *http.Request) {
	f, err := staticfs.StaticFiles.ReadFile("cmd/web/assets/index.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Write(f)
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

func GetCopyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("copied to clipboard"))
}
