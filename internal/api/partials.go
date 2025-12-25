package api

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

/*
type WidgetPage struct {
    Widgets []pages.Widget
}
*/

func partialsRoutes(serveMux *chi.Mux, s *Server) {
	mux := chi.NewRouter()

	mux.Get("/exercises", s.partialExercises)
	mux.Get("/workouts", s.partialWorkouts)

	serveMux.Mount("/partials", mux)
}

// todo: break this out
func postRoutes(serveMux *chi.Mux, s *Server) {
	serveMux.Post("/exercises", s.addExercise)
	serveMux.Patch("/exercises/{id}", s.patchExercise)
	serveMux.Delete("/exercises/{id}", s.deleteExercise)
}

func AddEmptyClick(name string) string {
	return fmt.Sprintf("onclick=\"%s()\"", name)
}

func DeleteClick(name string, id int) string {
	return fmt.Sprintf("onclick=\"%s(%d)\"", name, id)
}

func EditClick(name string, id int) string {
	return fmt.Sprintf("onclick=\"%s(%d)\"", name, id)
}

func EditCancel(name string, id int) string {
	return fmt.Sprintf("onclick=\"%s(%d)\"", name, id)
}

func PatchClick(name string, id int) string {
	return fmt.Sprintf("onclick=\"%s(%d)\"", name, id)
}

func (s *Server) partialExercises(w http.ResponseWriter, r *http.Request) {
	e := s.exercise.GetAll()

	funcMap := template.FuncMap{
		"deleteClick": func(id int) template.HTMLAttr {
			return template.HTMLAttr(DeleteClick("deleteExercise", id))
		},
		"editClick": func(id int) template.HTMLAttr {
			return template.HTMLAttr(EditClick("editExercise", id))
		},
		"editCancel": func(id int) template.HTMLAttr {
			return template.HTMLAttr(EditCancel("editExerciseCancel", id))
		},
		"patchClick": func(id int) template.HTMLAttr {
			return template.HTMLAttr(PatchClick("patchListener", id))
		},
	}

	err := s.tpls.Funcs(funcMap).ExecuteTemplate(w, "partials/exercise-list.html", e)
	if err != nil {
		fmt.Println("couldn't open exercise list partial", err)
	}
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

func (s *Server) partialWorkouts(w http.ResponseWriter, r *http.Request) {
	wo := s.workout.GetAll()

	funcMap := template.FuncMap{
		"addClick": func() template.HTMLAttr {
            return template.HTMLAttr(AddEmptyClick("addWorkout"))
        },
	}

	err := s.tpls.Funcs(funcMap).ExecuteTemplate(w, "partials/workout-list.html", wo)
	if err != nil {
		fmt.Println("couldn't open exercise list partial", err)
	}
}

/*
func (s *Server) partialsWidgets(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")

    widgets := pages.GetWidgets(idStr)
    widgetPage := WidgetPage{Widgets: widgets}

    err :=  s.tpls.ExecuteTemplate(w, "partials/widget-list.html", widgetPage)
    if err != nil {
        fmt.Println("couldn't open widget partials", err)
    }
}
*/
