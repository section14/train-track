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
	serveMux.Get("/workouts", s.workoutPage)
}

func appendPath(p string) string {
    return fmt.Sprintf("/static/js/extracted%s", p)
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/workouts", http.StatusFound)
}

func (s *Server) exercisePage(w http.ResponseWriter, r *http.Request) {
	head := HeadImports{
		Title: "Exercises",
		Js: []string{
            // Js to import for this page
			appendPath("/pages/exercises.js"),
			appendPath("/partials/exercise-list.js"),
		},
	}

	err := s.tpls.ExecuteTemplate(w, "pages/exercises.html", head)
	if err != nil {
		fmt.Println("couldn't open exercises page", err)
	}
}

func (s *Server) workoutPage(w http.ResponseWriter, r *http.Request) {
	head := HeadImports{
		Title: "Workouts",
		Js:    []string{
            // Js to import for this page
            appendPath("/pages/workouts.js"),
			appendPath("/pages/exercises.js"),
			appendPath("/partials/workout-list.js"),
			appendPath("/partials/workout-item.js"),
			appendPath("/partials/exercise-list.js"),
        },
	}

	err := s.tpls.ExecuteTemplate(w, "pages/workouts.html", head)
	if err != nil {
		fmt.Println("couldn't open workouts page", err)
	}
}
