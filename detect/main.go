package main

import (
	"os"
)

func main() {
	if len(os.Args) != 4 {
		panic("Expected 4 args!")
	}
	// os.Args[0] is just name of program
	platformVarsPath := os.Args[1]
	planPath := os.Args[2]
	appPath := os.Args[3]
	returnStatus, err := DetectionFunction(platformVarsPath, planPath, appPath)
	if err != nil {
		panic(err)
	}
	os.Exit(returnStatus)
}
