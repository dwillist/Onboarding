package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-community/mri"
)

func main() {
	packageJSONParser := mri.NewPackageJSONParser()
	entryResolver := mri.NewPlanEntryResolver()
	dependencyManager := postal.NewService(cargo.NewTransport())

	packit.Run(
		mri.Detect(packageJSONParser),
		mri.Build(
			entryResolver,
			dependencyManager,
		),
	)
}
