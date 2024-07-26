package server

import (
	"io/fs"
	"log"
	"net/http"

	embedfs "github.com/mloughton/pbifiltergen/cmd/web"
	"github.com/mloughton/pbifiltergen/internal/dax"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", GetHealthHandler)

	mux.Handle("/", AppHandler())

	mux.HandleFunc("POST /input", s.limit(PostInputHandler))

	return mux
}

func AppHandler() http.Handler {
	fs, err := fs.Sub(embedfs.StaticFiles, "assets")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fs))
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
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	dax, err := dax.GenerateDax(columns)
	if err != nil {
		w.Header().Add("HX-Retarget", "#error")
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	_, err = w.Write([]byte(dax))
	if err != nil {
		log.Fatal(err)
	}
}
