package mri

import (
	"encoding/json"
	"fmt"
	"os"
)

type PackageJSONParser struct{}

func NewPackageJSONParser() PackageJSONParser {
	return PackageJSONParser{}
}

type PackageJSON struct {
	Name         string `json:"name"`
	Dependencies struct {
		MRI string `json:"MRI"`
	} `json:"dependencies"`
}

func (p PackageJSONParser) ParseVersion(path string) (string, error) {
	var packageJSON PackageJSON

	file, err := os.Open(path)
	switch {
	case os.IsNotExist(err):
		return "", nil
	default:
		return "", fmt.Errorf("error opening package.json file: %s", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(packageJSON)
	if err != nil {
		return "", fmt.Errorf("unable to parse package.json file")
	}

	return packageJSON.Dependencies.MRI, nil
}
