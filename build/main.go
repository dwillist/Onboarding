package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 5 {
		panic("Build expects 5 args!")
	}
	// os.Args[0] is just name of program
	buildpackPath := filepath.Join(filepath.Dir(os.Args[0]), "buildpack.toml")
	layersPath := os.Args[1]
	platformPath := os.Args[2]
	planPath := os.Args[3]
	appPath := os.Args[4]

	builder := NewBuilder(&http.Client{})
	returnStatus, err := builder.BuildFunction(buildpackPath, layersPath, platformPath, planPath, appPath)
	if err != nil {
		panic(err)
	}
	os.Exit(returnStatus)
}
