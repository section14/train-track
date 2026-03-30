package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func partialsRoutes(serveMux *chi.Mux, s *Server) {
	mux := chi.NewRouter()

	mux.Get("/exercises", s.partialExercises)
	mux.Get("/workouts", s.partialWorkouts)
	mux.Get("/workouts/{id}", s.partialWorkout)

	serveMux.Mount("/partials", mux)
}

func (s *Server) partialExercises(w http.ResponseWriter, r *http.Request) {
	e := s.exercise.GetAll()

	err := s.tpls.ExecuteTemplate(w, "partials/exercise-list.html", e)
	if err != nil {
		fmt.Println("couldn't open exercise list partial", err)
	}
}

func (s *Server) partialWorkouts(w http.ResponseWriter, r *http.Request) {
	wo := s.workout.GetAll()

	err := s.tpls.ExecuteTemplate(w, "partials/workout-list.html", wo)
	if err != nil {
		fmt.Println("couldn't open workout list partial", err)
	}
}

func (s *Server) partialWorkout(w http.ResponseWriter, r *http.Request) {
	wo, err := s.workout.Get(r)
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	err = s.tpls.ExecuteTemplate(w, "partials/workout-item.html", wo)
	if err != nil {
		fmt.Println("couldn't open workout item partial", err)
	}
}
