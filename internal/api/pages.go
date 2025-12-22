package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HeadImports struct {
	Title string
	Js    []string
}

func pageRoutes(serveMux *chi.Mux, s *Server) {
	serveMux.Get("/", s.homePage)
	serveMux.Get("/exercises", s.exercisePage)
}

func appendPath(p string) string {
    return fmt.Sprintf("/static/js/extracted%s", p)
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	head := HeadImports{
		Title: "Home",
		Js:    nil,
	}

	err := s.tpls.ExecuteTemplate(w, "pages/home.html", head)
	if err != nil {
		fmt.Println("couldn't open home", err)
	}
}

func (s *Server) exercisePage(w http.ResponseWriter, r *http.Request) {
	head := HeadImports{
		Title: "Exercises",
        //make a function that adds to beginning path automatically
		Js: []string{
			appendPath("/pages/exercises.js"),
			appendPath("/partials/exercise-list.js"),
		},
	}

	err := s.tpls.ExecuteTemplate(w, "pages/exercises.html", head)
	if err != nil {
		fmt.Println("couldn't open widgets", err)
	}
}
