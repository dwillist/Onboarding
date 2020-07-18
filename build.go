package mri

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/postal"
)

//go:generate faux --interface EntryResolver --output fakes/entry_resolver.go
type EntryResolver interface {
	Resolve([]packit.BuildpackPlanEntry) packit.BuildpackPlanEntry
}

//go:generate faux --interface DependencyManager --output fakes/dependency_manager.go
type DependencyManager interface {
	Resolve(path, id, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

//go:generate faux --interface BuildPlanRefinery --output fakes/build_plan_refinery.go
type BuildPlanRefinery interface {
	BillOfMaterial(dependency postal.Dependency) packit.BuildpackPlan
}

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) error
}

func Build(entries EntryResolver, dependencies DependencyManager) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		fmt.Println("Packit onboarding Buildpack!")
		fmt.Println("--Resolving dependency version")

		entry := entries.Resolve(context.Plan.Entries)

		mri, err := dependencies.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), entry.Name, entry.Version, context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		fmt.Printf("--Resolved dependency to: %s, %s\n", mri.Name, mri.Version)

		mriLayer, err := context.Layers.Get(MRI)
		if err != nil {
			return packit.BuildResult{}, err
		}

		mriLayer.Launch = entry.Metadata["launch"] == true
		mriLayer.Build = entry.Metadata["build"] == true
		mriLayer.Cache = entry.Metadata["build"] == true

		fmt.Println("--Executing build process")

		err = mriLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		fmt.Printf("--Installing MRI %s\n", mri.Version)
		err = dependencies.Install(mri, context.CNBPath, mriLayer.Path)
		if err != nil {
			return packit.BuildResult{}, err
		}

		fmt.Println("--Completed in")

		return packit.BuildResult{
			Plan:   context.Plan, // really we should update this to have more information in it
			Layers: []packit.Layer{mriLayer},
		}, nil
	}
}
