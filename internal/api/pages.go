package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func pageRoutes(serveMux *chi.Mux, s *Server) {
    serveMux.Get("/", s.homePage)
    serveMux.Get("/exercises", s.exercisePage)
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
    err := s.tpls.ExecuteTemplate(w, "pages/home.html", nil)
    if err != nil {
        fmt.Println("couldn't open home", err)
    }
}

func (s *Server) exercisePage(w http.ResponseWriter, r *http.Request) {
    err := s.tpls.ExecuteTemplate(w, "pages/exercises.html", nil)
    if err != nil {
        fmt.Println("couldn't open widgets", err)
    }
}
