package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "TestBuild", func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect  = NewWithT(t).Expect
			baseDir string
		)
		it.Before(func() {
			var err error
			baseDir, err = ioutil.TempDir("", "build-test")
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(os.RemoveAll(baseDir)).To(Succeed())
		})

		when("Build", func() {
			var (
				buildpackTOMLPath string
				layersPath        string
				platformPath      string
				planPath          string
				appPath           string
				builder           Builder
			)
			it.Before(func() {

				//
				// create a mock http server that will be used by the `downloadHelper` method
				// so we are not actually using a network
				// should be automatically used as long as we pass `downloadHelper` the correct uri
				//
				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, req *http.Request) {
							switch req.URL.Path {
							case "/some-download-url":
								tarReader, err := os.Open("testdata/test_tar.tgz")
								Expect(err).NotTo(HaveOccurred())

								_, err = io.Copy(w, tarReader)
								Expect(err).NotTo(HaveOccurred())
							default:
								http.NotFound(w, req)
							}
						},
					),
				)

				// initalize the builder that we are going to use to test the BuildFunction method
				builder = NewBuilder(server.Client())

				// write a buildpack.toml at the 'root' of our fake buildpack
				// this file follows the outline specified in the spec:
				// https://github.com/buildpacks/spec/blob/main/buildpack.md#buildpacktoml-toml
				buildpackTOMLPath = filepath.Join(baseDir, "buildpack.toml")
				Expect(ioutil.WriteFile(buildpackTOMLPath, []byte(fmt.Sprintf(`
[buildpack]
  id = "some-test/buildpack"
  name = "Super cool test"

[metadata]
  [[metadata.dependencies]]
    id = "some-dependency"
    sha256 = "NA"
    stacks = ["io.buildpacks.stacks.bionic"]
	uri = "%s/some-download-url"
    version = "14.5.0"

[[stacks]]
  id = "io.buildpacks.stacks.bionic"`, server.URL)), os.ModePerm)).To(Succeed())

				// path to the layers directory, this is where we should untar our 'node' dependency
				// more info in the spec:
				// https://github.com/buildpacks/spec/blob/main/buildpack.md#build
				// https://github.com/buildpacks/spec/blob/main/buildpack.md#layers
				layersPath = filepath.Join(baseDir, "layers")
				Expect(os.MkdirAll(layersPath, os.ModePerm)).To(Succeed())

				// unused in the solution, but required input to build
				platformPath = filepath.Join(baseDir, "platform")
				Expect(os.MkdirAll(platformPath, os.ModePerm)).To(Succeed())

				// the BuildPack Plan file that contains info about the dependency we are going to install
				// used to pass info from detect to build
				// un-used in solution
				planPath = filepath.Join(baseDir, "plan.toml")

				// Path to application,
				// again unused in the solution implementation
				appPath = filepath.Join(baseDir, "app-dir")
				Expect(os.MkdirAll(appPath, os.ModePerm)).To(Succeed())

			})

			when("Fresh Build", func() {
				it("creates a layer and installs node dependency", func() {
					returnVal, err := builder.BuildFunction(buildpackTOMLPath, layersPath, platformPath, planPath, appPath)
					Expect(err).NotTo(HaveOccurred())

					Expect(returnVal).To(Equal(0))

					// if downloadHelper is called with the correct destination, the below assertions will pass
					Expect(filepath.Join(layersPath, "node")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "file.txt")).To(BeAnExistingFile())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir", "inner_file.txt")).To(BeAnExistingFile())

					// Need to write this file, for this example its ok to just hard code the truth values assigned
					// to launch, build, and cache
					Expect(filepath.Join(layersPath, "node.toml")).To(BeAnExistingFile())
					nodeTOMLContents, err := ioutil.ReadFile(filepath.Join(layersPath, "node.toml"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(nodeTOMLContents)).To(Equal(`launch = true
build = false
cache = false
`))
				})
			})

			// Very similar test to the above, however `pack` will restore the node.toml file
			// if building a second time. For our simple implementation we are not going to do
			// any optimization and should over-write the contents of node.toml
			when("There are existing layer contents", func() {
				it.Before(func() {
					Expect(ioutil.WriteFile(filepath.Join(layersPath, "node.toml"), []byte(`inital contents`), os.ModePerm)).To(Succeed())
				})
				it("deletes them before creating layer and installing node dependency", func() {
					returnVal, err := builder.BuildFunction(buildpackTOMLPath, layersPath, platformPath, planPath, appPath)
					Expect(err).NotTo(HaveOccurred())

					Expect(returnVal).To(Equal(0))

					Expect(filepath.Join(layersPath, "node")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "file.txt")).To(BeAnExistingFile())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir", "inner_file.txt")).To(BeAnExistingFile())

					Expect(filepath.Join(layersPath, "node.toml")).To(BeAnExistingFile())
					nodeTOMLContents, err := ioutil.ReadFile(filepath.Join(layersPath, "node.toml"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(nodeTOMLContents)).To(Equal(`launch = true
build = false
cache = false
`))

				})
			})

		})
	})
}
