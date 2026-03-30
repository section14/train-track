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
    if len(os.Args) < 4 {
        fmt.Println("Not enough arguments. Please provide a run type, host and port. ex: prod localhost 8080")
        os.Exit(1)
    }

    buildType := os.Args[1]
    host := os.Args[2]
    port := os.Args[3]

    if buildType == "dev" {
		api.ServeDev(host, port)
	} else if buildType == "prod" {
		api.ServeDev(host, port)
        //todo: implement embedding
		//api.ServeProd(host, port, templates, static)
	} else {
		fmt.Printf("%s is not a valid run type. Supply dev or prod.", buildType)
		os.Exit(1)
	}
}
