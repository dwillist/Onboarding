package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 4 {
		panic("Build expects 4 args!")
	}
	buildpackPath, err := filepath.Abs(filepath.Join(os.Args[0], "..", "..", "buildpack.toml"))
	if err != nil {
		panic(fmt.Sprintf("error getting path to buildpack.toml: %s", err))
	}
	layersPath := os.Args[1]
	platformPath := os.Args[2]
	planPath := os.Args[3]
	appPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("error getting the working directory: %s", err))
	}

	builder := NewBuilder(&http.Client{})
	returnStatus, err := builder.BuildFunction(buildpackPath, layersPath, platformPath, planPath, appPath)
	if err != nil {
		panic(err)
	}
	os.Exit(returnStatus)
}
