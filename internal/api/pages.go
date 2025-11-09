package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func pageRoutes(mux *chi.Mux, s *Server) {
    mux.Get("/", s.homePage)
    mux.Get("/widgets", s.widgetPage)
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
    err := s.tpls.ExecuteTemplate(w, "pages/home.html", nil)
    if err != nil {
        fmt.Println("couldn't open home", err)
    }
}

func (s *Server) widgetPage(w http.ResponseWriter, r *http.Request) {
    err := s.tpls.ExecuteTemplate(w, "pages/widgets.html", nil)
    if err != nil {
        fmt.Println("couldn't open widgets", err)
    }
}
