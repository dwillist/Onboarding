package mri_test

import (
	"testing"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-community/mri"
	"github.com/paketo-community/mri/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		packageJSONParser *fakes.VersionParser
		detect            packit.DetectFunc
	)

	it.Before(func() {
		packageJSONParser = &fakes.VersionParser{}

		detect = mri.Detect(packageJSONParser)
	})

	it("returns a plan that provides mri", func() {
		result, err := detect(packit.DetectContext{
			WorkingDir: "/working-dir",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Plan).To(Equal(packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{
				{Name: mri.MRI},
			},
		}))
	})

	context("when the source code contains a package.json file", func() {
		it.Before(func() {
			packageJSONParser.ParseVersionCall.Returns.Version = "4.5.6"
		})

		it("returns a plan that provides and requires that version of mri", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: "/working-dir",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: mri.MRI},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name:    mri.MRI,
						Version: "4.5.6",
					},
				},
			}))

			Expect(packageJSONParser.ParseVersionCall.Receives.Path).To(Equal("/working-dir/package.json"))
		})
	})
}
