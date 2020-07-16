package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "TestDetect", func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect
		)

		when("Detect", func() {
			var (
				platformPath string
				planPath     string
				appPath      string
				detector     Detector
			)

			// Initializing some test directories and files that are needed by the
			// DetectFunction call

			it.Before(func() {
				baseDir, err := ioutil.TempDir("", "testDir")
				Expect(err).NotTo(HaveOccurred())

				// path to the Build Plan file, this file is used to
				// specify what this buildpack will be providing, so that subsequent
				// buildpacks know if they can be used,
				// Ex:
				// in order to run 'npm install', we need 'node' & 'npm'
				// available
				//
				// Specification for what this file looks like here:
				// https://github.com/buildpacks/spec/blob/main/buildpack.md#build-plan-toml
				planPath = filepath.Join(baseDir, "plan.toml")

				// unused in the solution, but required input to build
				platformPath = filepath.Join(baseDir, "platform")
				Expect(os.MkdirAll(platformPath, os.ModePerm)).To(Succeed())

				// Path to application,
				// needed to figure out if we have the correct 'name' entry in our
				// package.json
				appPath = filepath.Join(baseDir, "application")
				Expect(os.MkdirAll(appPath, os.ModePerm)).To(Succeed())

				// initialize the detector to test against
				detector = NewDetector()
			})

			// When an application does not meet our detection criteria
			// we should return with an exit status of 100
			// as outlined in the spec:
			// https://github.com/buildpacks/spec/blob/main/buildpack.md#detection
			when("when application has no package.json", func() {
				it("fails detection", func() {
					exitStatus, err := detector.DetectFunction(platformPath, planPath, appPath)
					Expect(exitStatus).To(Equal(100))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			when(`when application has "name": "onboarding_app" in package.json`, func() {
				it.Before(func() {
					err := ioutil.WriteFile(filepath.Join(appPath, "package.json"), []byte(`
{
  "name": "onboarding_app",
  "engines": { "node": "some-semver-version" }
}
`), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				})
				it("passes detection", func() {
					exitStatus, err := detector.DetectFunction(platformPath, planPath, appPath)
					Expect(exitStatus).To(Equal(0))
					Expect(err).NotTo(HaveOccurred())

					contents, err := ioutil.ReadFile(planPath)
					Expect(err).NotTo(HaveOccurred())
					fmt.Println(string(contents))
					Expect(string(contents)).To(Equal(`[[provides]]
  name = "node"

[[require]]
  name = "node"
  version = "some-semver-version"
`))
				})
			})

			when(`when application has other "name" in  package.json`, func() {
				it.Before(func() {
					err := ioutil.WriteFile(filepath.Join(appPath, "package.json"), []byte(`{"name": "blerb"}`), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				})
				it("fails detection", func() {
					exitStatus, err := detector.DetectFunction(platformPath, planPath, appPath)
					Expect(exitStatus).To(Equal(100))
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
}
