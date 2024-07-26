package server

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(s.limiter.Tokens())
		if !s.limiter.Allow() {
			w.Header().Add("HX-Retarget", "#error")
			_, err := w.Write([]byte("too many requests"))
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}
