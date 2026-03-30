package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func jsonRoutes(serveMux *chi.Mux, s *Server) {
	serveMux.Get("/json/exercises", s.exerciseJson)
	serveMux.Post("/exercises", s.addExercise)
	serveMux.Patch("/exercises/{id}", s.patchExercise)
	serveMux.Delete("/exercises/{id}", s.deleteExercise)

	serveMux.Post("/workouts", s.addWorkout)
	serveMux.Delete("/workouts/{id}", s.deleteWorkout)
	serveMux.Patch("/workouts/add/movement/{id}", s.addMovement)
	serveMux.Patch("/workouts/{workoutId}/{movementId}", s.updateWorkout)
    serveMux.Delete("/workouts/{workoutId}/{movementId}", s.deleteMovement)
}

func (s *Server) exerciseJson(w http.ResponseWriter, r *http.Request) {
	e := s.exercise.GetAllJson()
	w.Write(e)
}

func (s *Server) addExercise(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	err := s.exercise.Add(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) patchExercise(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	name := r.FormValue("name")
	err = s.exercise.Update(int(id), name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) deleteExercise(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = s.exercise.Delete(int(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) addWorkout(w http.ResponseWriter, r *http.Request) {
	err := s.workout.Add()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) deleteWorkout(w http.ResponseWriter, r *http.Request) {
	err := s.workout.Delete(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) addMovement(w http.ResponseWriter, r *http.Request) {
    fmt.Println("called")

	wo, err := s.workout.AddMovement(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = s.tpls.ExecuteTemplate(w, "partials/workout-item.html", wo)
	if err != nil {
		fmt.Println("couldn't open workout item partial", err)
	}
}

func (s *Server) deleteMovement(w http.ResponseWriter, r *http.Request) {
	wo, err := s.workout.DeleteMovement(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = s.tpls.ExecuteTemplate(w, "partials/workout-item.html", wo)
	if err != nil {
		fmt.Println("couldn't open workout item partial", err)
	}
}

func (s *Server) updateWorkout(w http.ResponseWriter, r *http.Request) {
	wo, err := s.workout.Update(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = s.tpls.ExecuteTemplate(w, "partials/workout-list.html", wo)
	if err != nil {
		fmt.Println("couldn't open workout list partial", err)
	}
}


