package api

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/section14/train-track/internal/config"
	"github.com/section14/train-track/internal/service"
	"github.com/section14/train-track/internal/store"
)

type Server struct {
	tpls *template.Template
    exercise *service.ExerciseService
    workout *service.WorkoutService
}

func NewServer(
    t *template.Template, 
    exercise *service.ExerciseService, 
    workout *service.WorkoutService) *Server {

	return &Server{
        tpls: t,
        exercise: exercise,
        workout: workout,
    }
}

func handlers(mux *chi.Mux, s *Server) {
	pageRoutes(mux, s)

    apiMux := chi.NewRouter()
    partialsRoutes(apiMux, s)
    postRoutes(apiMux, s)

    mux.Mount("/api", apiMux)
}

func systemTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
    cleanRoot := filepath.Clean(rootDir)
    pfx := len(cleanRoot)+1
    root := template.New("")

    //temp func test
    //todo: I think you're going to have to write forward "dummy" declarations like this
    root.Funcs(template.FuncMap{
        "clicker": func(id int) template.HTMLAttr {return ""},
    })

    err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
        if !info.IsDir() && strings.HasSuffix(path, ".html") {
            if e1 != nil {
                return e1
            }

            b, e2 := os.ReadFile(path)
            if e2 != nil {
                return e2
            }

            name := path[pfx:]
            t := root.New(name).Funcs(funcMap)
            _, e2 = t.Parse(string(b))
            if e2 != nil {
                return e2
            }
        }

        return nil
    })

    return root, err
}

func embeddedTemplates(files fs.FS, rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	//pfx := len(cleanRoot) + 1
	root := template.New("")

	err := fs.WalkDir(files, cleanRoot, func(path string, d fs.DirEntry, e1 error) error {
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			//b, e2 := os.ReadFile(path)
			b, e2 := fs.ReadFile(files, path)
			if e2 != nil {
				return e2
			}

			//name := path[pfx:]
			name := path

			t := root.New(name).Funcs(funcMap)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}

		}

		return nil
	})

	return root, err
}

func ServeDev() {
	mux := chi.NewRouter()
	//mux.Use(middleware.Logger)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

    currentDir,_ := os.Getwd()

    t, err := systemTemplates(fmt.Sprintf("%s/templates", currentDir), nil)
	if err != nil {
		log.Fatal("couldn't parse templates: ", err)
	}

	//static files
    staticFS := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	mux.Handle("/static/*", staticFS)

    fmt.Println("serving dev...")

	Serve(mux, t)
}

func ServeProd(templates embed.FS, static embed.FS) {
	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		log.Fatal("couldn't setup embedded templates directory: ", err)
	}

	//todo: you don't need this I guess?
	staticSub, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal("couldn't setup static sub directory: ", err)
	}

	fs.WalkDir(staticSub, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println("static?: ", path)
		return nil
	})

	//t, err := template.ParseFS(tpl)
	t, err := embeddedTemplates(sub, "", nil)
	if err != nil {
		log.Fatal("couldn't parse templates: ", err)
	}

	mux := chi.NewRouter()
	//mux.Use(middleware.Logger)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//static files
	mux.Handle("/static/*", http.FileServerFS(static))

	Serve(mux, t)
}

func Serve(mux *chi.Mux, t *template.Template) {
    env := config.NewEnv()

    //stores
    exerciseStore := store.NewExerciseStore(env)
    workoutStore := store.NewWorkoutStore(env)

    //services
    exerciseService := service.NewExerciseService(exerciseStore)
    workoutService := service.NewWorkoutService(workoutStore)

	server := NewServer(t, exerciseService, workoutService)
	handlers(mux, server)

	//addr := fmt.Sprintf("%s:%s", env.Location, env.Port)
	addr := fmt.Sprintf("%s:%s", "localhost", "8080")

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
		// other settings omitted
	}

	fmt.Println("serving on localhost:8080...")
	log.Fatal(s.ListenAndServe())
}
