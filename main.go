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

    if buildType == "dev" {
		api.ServeDev()
	} else if buildType == "prod" {
		api.ServeProd(templates, static)
	} else {
		fmt.Printf("%s is not a valid build type. Supply dev or prod.", buildType)
		os.Exit(1)
	}
}
