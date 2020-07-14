package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const AppName string = "onboarding_app"

//
// Inputs:
//   platformVarsPath: path to a folder that contains environment variables set by the platform,
//     the platform is the 'pack' tool in our case
//   planPath: this is the path to a .toml file that encodes the requirements that this buildpack will need
//     during it's build phase
//   appPath: path to the root of our application
//
func DetectionFunction(platformVarsPath, planPath, appPath string) (int, error) {
	var packageJSON struct {
		Name string `json:"name"`
	}
	packageJSONFile, err := os.Open(filepath.Join(appPath, "package.json"))
	switch {
	case os.IsNotExist(err):
		return 100, nil
	case err != nil:
		return -1, fmt.Errorf("error opening app's package.json file: %s", err)
	}

	err = json.NewDecoder(packageJSONFile).Decode(&packageJSON)
	if err != nil {
		return -1, fmt.Errorf("error decoding package.json file: %s", err)
	}

	if packageJSON.Name == AppName {
		return 0, nil
	}
	return 100, nil
}
