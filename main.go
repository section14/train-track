package main

import (
    "embed"
    "fmt"
    "os"

    "github.com/section14/train-track/internal/api"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

func main() {
    buildType := os.Args[1]
    port := os.Args[2]

    if len(os.Args) < 2 {
        fmt.Println("provide a run type and port. ex: prod 8080")
    }

    if buildType == "dev" {
		api.ServeDev(port)
	} else if buildType == "prod" {
		api.ServeProd(port, templates, static)
	} else {
		fmt.Printf("%s is not a valid run type. Supply dev or prod.", buildType)
		os.Exit(1)
	}
}
