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
    mux.Post("/exercises", s.addExercise)

    serveMux.Mount("/partials", mux)
}

//todo: break this out
func postRoutes(serveMux *chi.Mux, s *Server) {
    serveMux.Post("/exercises", s.addExercise)
    serveMux.Delete("/exercises/{id}", s.deleteExercise)
}

func DeleteClick(id int) string {
   return fmt.Sprintf("onclick=\"deleteExercise(%d)\"", id)
}

func (s *Server) partialExercises(w http.ResponseWriter, r *http.Request) {
    e := s.exercise.GetAll()

    funcMap := template.FuncMap{
        "clicker": func(id int) template.HTMLAttr { return template.HTMLAttr(DeleteClick(id))},
    }

    err :=  s.tpls.Funcs(funcMap).ExecuteTemplate(w, "partials/exercise-list.html", e)
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
