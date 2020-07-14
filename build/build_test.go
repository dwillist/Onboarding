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

		//it.After(func() {
		//	Expect(os.RemoveAll(baseDir)).To(Succeed())
		//})

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

				builder = NewBuilder(server.Client())

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

				layersPath = filepath.Join(baseDir, "layers")
				Expect(os.MkdirAll(layersPath, os.ModePerm)).To(Succeed())

				platformPath = filepath.Join(baseDir, "platform")
				Expect(os.MkdirAll(platformPath, os.ModePerm)).To(Succeed())

				planPath = filepath.Join(baseDir, "plan.toml")

				appPath = filepath.Join(baseDir, "app-dir")
				Expect(os.MkdirAll(appPath, os.ModePerm)).To(Succeed())

			})

			when("First Build", func() {
				it("creates a layer and installs node dependency", func() {
					returnVal, err := builder.BuildFunction(buildpackTOMLPath, layersPath, platformPath, planPath, appPath)
					Expect(err).NotTo(HaveOccurred())

					Expect(returnVal).To(Equal(0))

					Expect(filepath.Join(layersPath, "node")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node.toml")).To(BeAnExistingFile())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "file.txt")).To(BeAnExistingFile())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir", "inner_file.txt")).To(BeAnExistingFile())

					nodeTOMLContents, err := ioutil.ReadFile(filepath.Join(layersPath, "node.toml"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(nodeTOMLContents)).To(Equal(`launch = "true"
build = "true"
cache = "true"
`))
				})
			})

			when("There are existing layer contents", func() {
				it.Before(func() {
					Expect(ioutil.WriteFile(filepath.Join(layersPath, "node.toml"), []byte(`inital contents`), os.ModePerm)).To(Succeed())
				})
				it("deletes them before creating layer and installing node dependency", func() {
					returnVal, err := builder.BuildFunction(buildpackTOMLPath, layersPath, platformPath, planPath, appPath)
					Expect(err).NotTo(HaveOccurred())

					Expect(returnVal).To(Equal(0))

					Expect(filepath.Join(layersPath, "node")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node.toml")).To(BeAnExistingFile())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "file.txt")).To(BeAnExistingFile())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir")).To(BeADirectory())
					Expect(filepath.Join(layersPath, "node", "fake_archive_root", "inner_dir", "inner_file.txt")).To(BeAnExistingFile())

					nodeTOMLContents, err := ioutil.ReadFile(filepath.Join(layersPath, "node.toml"))
					Expect(err).NotTo(HaveOccurred())

					Expect(string(nodeTOMLContents)).To(Equal(`launch = "true"
build = "true"
cache = "true"
`))

				})
			})

		})
	})
}
