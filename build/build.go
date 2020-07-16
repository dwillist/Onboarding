package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/blang/semver"
)

// golang struct deffinition of buildpack.toml
// Buildpack.toml specification here:
// https://github.com/buildpacks/spec/blob/main/buildpack.md#buildpacktoml-toml
// Notice we ignore the 'order' section entirely, as is irrelevant for
// implementation buildpacks
type BuildpackTOMLStruct struct {
	BuildpackFields struct {
		id   string `toml:"id"`
		name string `toml:name"`
	} `toml:"buildpack"`
	Metadata struct {
		Dependencies []struct {
			ID      string   `toml:"id"`
			SHA256  string   `toml:sha256"`
			Stacks  []string `toml:"stacks"`
			URI     string   `toml:"uri"`
			Version string   `toml:"version"`
		} `toml:"dependencies"`
	} `toml:"metadata"`
}

type BuildPlanTOMLStruct struct {
	Entries []struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
	} `toml:"entries"`
	// Also metadata fields that we are going to ignore for now
}

// golang struct deffinition of the <layer>.toml struct
// <layer>.toml specification here:
// https://github.com/buildpacks/spec/blob/main/buildpack.md#layer-content-metadata-toml
type LayerTOMLStruct struct {
	Launch bool `toml:"launch"`
	Build  bool `toml:"build"`
	Cache  bool `toml:"cache"`
	// Also metadata fields that we are going to ignore for now
}

type Builder struct {
	Client *http.Client
}

//
// Constructor that takes a http.Client
// needed to mock out requests in tests
//
func NewBuilder(client *http.Client) Builder {
	return Builder{
		Client: client,
	}
}

func (b Builder) BuildFunction(buildpackTOMLPath, layersDir, platformDir, planPath, appDir string) (int, error) {
	nodeLayerPath := filepath.Join(layersDir, "node")
	nodeLayerTOML := filepath.Join(layersDir, "node.toml")

	var (
		buildpackTOMLStruct BuildpackTOMLStruct
		buildPlanTOMLStruct BuildPlanTOMLStruct
	)

	//
	// Decode the contents of the buildpack.toml file,
	// we need the uri to the node dependency we want to download
	//
	fmt.Println("--- Decoding buildpack.toml file")
	buildpackTOMLContents, err := ioutil.ReadFile(buildpackTOMLPath)
	if err != nil {
		return 100, fmt.Errorf("unable to read buildpack.toml file: %s", err)
	}

	//
	// Decode the the buildpack.toml file into the BuildpackTOMLStruct
	// defined above using the toml library
	//
	_, err = toml.Decode(string(buildpackTOMLContents), &buildpackTOMLStruct)
	if err != nil {
		return 100, fmt.Errorf("unable to decode buildpack.toml file: %s", err)
	}

	//
	// just a quick sanity check
	// there is only one value
	//
	if len(buildpackTOMLStruct.Metadata.Dependencies) != 1 {
		return 100, fmt.Errorf("unexpected number of dependencies for our fake buildpack")
	}

	buildPlanContents, err := ioutil.ReadFile(planPath)
	if err != nil {
		return 100, fmt.Errorf("failed to read the BuildPlan toml file: %s", err)
	}

	_, err = toml.Decode(string(buildPlanContents), &buildPlanTOMLStruct)
	if err != nil {
		return 100, fmt.Errorf("failed to decode the BuildPlan toml file: %s", err)
	}

	//
	// Check to match the version given by the buildplan
	// with the actual dependency version that is in our buildpack.toml
	// to make sure that they agree
	//
	match, err := semverMatch(buildPlanTOMLStruct, buildpackTOMLStruct.Metadata.Dependencies[0].Version)
	if err != nil {
		return 100, fmt.Errorf("error matching buildplan and buildpack.toml 'node' versions: %s", err)
	}
	if !match {
		return 100, errors.New("no match for version constraint in buildpack.toml")
	}

	//
	// Use the download helper to, grab the node.tgz file
	// and unzip it into the nodeLayer
	//
	fmt.Println("--- Downloading node dependnecy")
	err = b.downloadHelper(buildpackTOMLStruct.Metadata.Dependencies[0].URI, nodeLayerPath)
	if err != nil {
		return 100, fmt.Errorf("unable to download node artifact: %s", err)
	}

	//
	// now write the node.toml file, forcing the launch
	// flag to be true so it ends up in our final image
	//
	fmt.Println("--- Writing node.toml file")
	nodeLayerFile, err := os.OpenFile(nodeLayerTOML, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return 100, fmt.Errorf("unable to open node.toml file for writing: %s", err)
	}
	defer nodeLayerFile.Close()

	// we are not going to worry too much about the flag values here
	// they will come up in subsequent examples
	nodeLayerTOMLContents := LayerTOMLStruct{
		Launch: true,
		Build:  false,
		Cache:  false,
	}

	//
	// use the toml library to encode the node.toml structure
	//out as a file that sits directly next to the nodeLayer
	//
	err = toml.NewEncoder(nodeLayerFile).Encode(nodeLayerTOMLContents)
	if err != nil {
		return 100, errors.New("unabel to write node_layer.toml contents")
	}

	//
	// return 0 for a successful exit status!
	//
	fmt.Println("--- Success!")
	return 0, nil
}

func semverMatch(buildPlan BuildPlanTOMLStruct, buildpackTOMLVersion string) (bool, error) {
	var (
		versionConstraint semver.Range
		err               error
	)

	foundConstraint := false
	for _, entry := range buildPlan.Entries {
		if entry.Name == "node" {
			foundConstraint = true
			versionConstraint, err = semver.ParseRange(entry.Version)
			if err != nil {
				return false, fmt.Errorf("invalid version from BuildPlan: %s", err)
			}
		}
	}

	if !foundConstraint {
		return false, errors.New("Unable to find proper version constraint")
	}

	buildpackDependencyVersion, err := semver.Parse(buildpackTOMLVersion)
	if err != nil {
		return false, fmt.Errorf(
			"unable to parse buildpack.toml version %s: %s",
			buildpackTOMLVersion,
			err,
		)
	}

	return versionConstraint(buildpackDependencyVersion), nil
}

// downloads & unzips to a destination (we know we are useing .tar files here)
func (b Builder) downloadHelper(uri, dest string) error {
	resp, err := b.Client.Get(uri)
	if err != nil {
		return fmt.Errorf("fetching uri failed: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response status: %s", resp.Status)
	}

	defer resp.Body.Close()
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to create gzip reader on response: %s", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	if err != nil {
		return fmt.Errorf("unable to create tar reader on gzipReader: %s", err)
	}

	for {
		hdr, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return errors.New("error when reading .tar file")

		case hdr.FileInfo().IsDir():
			os.MkdirAll(filepath.Join(dest, hdr.Name), os.ModePerm)
		default: // assume we have a regular file, not 'production' ready code
			destFile, err := os.OpenFile(filepath.Join(dest, hdr.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return fmt.Errorf("unable to open dest path for file writing: %s", err)
			}
			_, err = io.Copy(destFile, tarReader)
			if err != nil {
				return fmt.Errorf("unable to copy tar file to dest path: %s", err)
			}
			destFile.Close()
		}
	}
}
