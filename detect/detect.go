package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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

type BuildPlan struct {
	Provides []Provide `toml:"provides"`
	Requires []Require `toml:"require"`
}

type Provide struct {
	Name string `toml:"name"`
}

type Require struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type Detector struct{}

func NewDetector() Detector {
	return Detector{}
}

func (d Detector) DetectFunction(platformVarsPath, planPath, appPath string) (int, error) {
	var packageJSON struct {
		Name    string `json:"name"`
		Engines struct {
			NodeVersion string `json:"node"`
		} `json:"engines"`
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

	if packageJSON.Name != AppName {
		return 100, nil
	}

	// Write out our Buildplan
	planFile, err := os.OpenFile(planPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return -1, fmt.Errorf("error opening planPath for writing: %s", err)
	}

	buildPlan := BuildPlan{
		Provides: []Provide{
			{
				Name: "node",
			},
		},
		Requires: []Require{
			{
				Name:    "node",
				Version: packageJSON.Engines.NodeVersion,
			},
		},
	}

	err = toml.NewEncoder(planFile).Encode(&buildPlan)
	if err != nil {
		return -1, fmt.Errorf("error writing BuildPlan to toml file at planPath: %s", err)
	}

	return 0, nil
}
