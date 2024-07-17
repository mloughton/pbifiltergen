package server

import (
	"io/fs"
	"log"
	"net/http"

	embedfs "github.com/mloughton/pbifiltergen/cmd/web"
	"github.com/mloughton/pbifiltergen/internal/dax"
)

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", GetHealthHandler)

	mux.Handle("/", AppHandlerTest())

	mux.HandleFunc("POST /input", PostInputHandler)

	return mux
}

func AppHandlerTest() http.Handler {
	fs, err := fs.Sub(embedfs.StaticFiles, "assets")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fs))
}
func AppHandler(w http.ResponseWriter, r *http.Request) {
	f, err := embedfs.StaticFiles.ReadFile("assets/index.html")
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
		return
	}
	dax, err := dax.GenerateDax(columns)
	if err != nil {
		w.Header().Add("HX-Retarget", "#error")
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(dax))
}
