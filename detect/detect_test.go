package main

import (
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
			it.Before(func() {
				baseDir, err := ioutil.TempDir("", "testDir")
				Expect(err).NotTo(HaveOccurred())
				planPath = filepath.Join(baseDir, "plan.toml")

				platformPath = filepath.Join(baseDir, "platform")
				Expect(os.MkdirAll(platformPath, os.ModePerm)).To(Succeed())

				appPath = filepath.Join(baseDir, "application")
				Expect(os.MkdirAll(appPath, os.ModePerm)).To(Succeed())

				detector = NewDetector()
			})

			when("when application has no package.json", func() {
				it("fails detection", func() {
					exitStatus, err := detector.DetectFunction(platformPath, planPath, appPath)
					Expect(exitStatus).To(Equal(100))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			when(`when application has "name": "onboarding_app" in package.json`, func() {
				it.Before(func() {
					err := ioutil.WriteFile(filepath.Join(appPath, "package.json"), []byte(`{"name": "onboarding_app"}`), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				})
				it("passes detection", func() {
					exitStatus, err := detector.DetectFunction(platformPath, planPath, appPath)
					Expect(exitStatus).To(Equal(0))
					Expect(err).NotTo(HaveOccurred())
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
