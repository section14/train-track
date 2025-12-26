package api

import (
	"embed"
	"errors"
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
	"github.com/section14/train-track/internal/extract"
	"github.com/section14/train-track/internal/service"
	"github.com/section14/train-track/internal/store"
)

type Server struct {
	tpls     *template.Template
	exercise *service.ExerciseService
	workout  *service.WorkoutService
}

func NewServer(
	t *template.Template,
	exercise *service.ExerciseService,
	workout *service.WorkoutService) *Server {

	return &Server{
		tpls:     t,
		exercise: exercise,
		workout:  workout,
	}
}

func handlers(mux *chi.Mux, s *Server) {
	pageRoutes(mux, s)

	apiMux := chi.NewRouter()
	partialsRoutes(apiMux, s)
	postRoutes(apiMux, s)

	mux.Mount("/api", apiMux)
}

type PagesJs struct {
	Name string
	Js   []string
}

func extractSystemTemplates(rootDir, targetDir, extractedDir string) (map[string]string, error) {
	cleanRoot := filepath.Clean(rootDir)
	//js := make([]string, 0)
	js := make(map[string]string)

	/*

	   Idea for importing JS:

	   Extract every page and partial to it's own js file, which mirrors the same folder/file
	   structure as templates. Then, supply a list to the ExecuteTemplate function in
	   api/pages.go, which will update the <head> with the correct imports

	*/

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			newNode, e2 := extract.ExtractJs(path)
			if e2 != nil {
				return e2
			}

			//get current working directory
			wd, _ := os.Getwd()

			parts := strings.Split(path, fmt.Sprintf("%s/%s", wd, targetDir))
			htmlPath := filepath.Join(wd, extractedDir, parts[1])

			//create new HTML file and write to it
			err := os.MkdirAll(filepath.Dir(htmlPath), 0770)
			if err != nil {
				return err
			}

			file, err := os.Create(htmlPath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = file.Write(newNode.NewHTML)
			if err != nil {
				return err
			}

			if newNode.JsRemoved {
				//build a new path for the js file
				split := strings.Split(path, fmt.Sprintf("%s/%s", wd, "templates"))
				trimmed := strings.TrimSuffix(split[1], ".html")
				jsPath := fmt.Sprintf("%s%s%s%s", wd, "/static/js/extracted", trimmed, ".js")

				//js = append(js, string(newNode.Js))
				js[jsPath] = string(newNode.Js)
			}

		}

		return nil
	})

	return js, err
}

func buildJsFile(currentDir string, data map[string]string) error {
	for jsFileName, jsFileData := range data {

		//create new JS file (plus optional directory) and write to it
		err := os.MkdirAll(filepath.Dir(jsFileName), 0770)
		if err != nil {
			return err
		}

		file, err := os.Create(jsFileName)
		if err != nil {
			return errors.New(fmt.Sprintf("couldn't open extracted.js file %s", err))
		}
		defer file.Close()

		//var sb strings.Builder

		//for _, j := range jsFileData {
		//sb.WriteString(jsFileData)
		//sb.WriteString("\n")
		//}
		_, err = file.WriteString(jsFileData)
		if err != nil {
			return errors.New(fmt.Sprintf("couldn't write to extracted.js %s", err))
		}
	}

	//build extracted js file
	/*
		file, err := os.Create(filepath.Join(currentDir, "static", "js", "extracted.js"))
		if err != nil {
			return errors.New(fmt.Sprintf("couldn't open extracted.js file %s", err))
		}
		defer file.Close()

		var sb strings.Builder

		for _, j := range data {
			sb.WriteString(j)
			sb.WriteString("\n")
		}

		_, err = file.WriteString(sb.String())
		if err != nil {
			return errors.New(fmt.Sprintf("couldn't write to extracted.js %s", err))
		}
	*/

	return nil
}

// todo: I think you're going to have to write forward "dummy" declarations like this
func forwardFuncs() template.FuncMap {
	return template.FuncMap{
		"addClick":    func(id int) template.HTMLAttr { return "" },
		"deleteClick": func(id int) template.HTMLAttr { return "" },
		"editClick":   func(id int) template.HTMLAttr { return "" },
		"editCancel":  func(id int) template.HTMLAttr { return "" },
		"patchClick":  func(id int) template.HTMLAttr { return "" },
		"selectClick": func(id int) template.HTMLAttr { return "" },
	}
}

func systemTemplates(
	root *template.Template,
	rootDir string,
	funcMap template.FuncMap) (*template.Template, error) {

	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	//root := template.New("")

	root.Funcs(forwardFuncs())

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

func embeddedTemplates(
	root *template.Template,
	files fs.FS,
	rootDir string,
	funcMap template.FuncMap) (*template.Template, error) {

	cleanRoot := filepath.Clean(rootDir)
	//pfx := len(cleanRoot) + 1
	//root := template.New("")

	root.Funcs(forwardFuncs())

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

	currentDir, _ := os.Getwd()
	templatesDir := fmt.Sprintf("%s/templates", currentDir)
	extractedDir := fmt.Sprintf("%s/extracted", currentDir)
	rootTemplates := template.New("")

	//parse web component <template> files
	wc, err := systemTemplates(rootTemplates, fmt.Sprintf("%s/static/components", currentDir), nil)
	if err != nil {
		log.Fatal("couldn't parse web component templates: ", err)
	}

	//extract
	js, err := extractSystemTemplates(templatesDir, "templates", "extracted")
	if err != nil {
		log.Fatal("couldn't extract templates: ", err)
	}

	err = buildJsFile(currentDir, js)
	if err != nil {
		log.Fatal(err)
	}

	//parse regular templates
	t, err := systemTemplates(wc, extractedDir, nil)
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
	rootTemplates := template.New("")

	t, err := embeddedTemplates(rootTemplates, sub, "", nil)
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
