package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		errorString := fmt.Sprint("Expected 3 args: instead got", strings.Join(os.Args, "\n"))
		panic(errorString)
	}
	// os.Args[0] is just name of program
	detector := NewDetector()
	platformVarsPath := os.Args[1]
	planPath := os.Args[2]
	appPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("error getting the working directory: %s", err))
	}
	returnStatus, err := detector.DetectFunction(platformVarsPath, planPath, appPath)
	if err != nil {
		panic(err)
	}
	os.Exit(returnStatus)
}
