package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

    "github.com/section14/evenflow/internal/ui/pages"
)

type WidgetPage struct {
    Widgets []pages.Widget
}

func partialsRoutes(mux *chi.Mux, s *Server) {
    pMux := chi.NewRouter()

    pMux.Get("/widgets/{id}", s.partialsWidgets)

    mux.Mount("/partials", pMux)
}

func (s *Server) partialsWidgets(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")

    widgets := pages.GetWidgets(idStr)
    widgetPage := WidgetPage{Widgets: widgets}

    err :=  s.tpls.ExecuteTemplate(w, "partials/widget-list.html", widgetPage)
    if err != nil {
        fmt.Println("couldn't open widget partials", err)
    }
}
