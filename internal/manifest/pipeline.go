// Package manifest implements a standard set of osbuild pipelines. A pipeline
// conceptually represents a named filesystem tree, optionally generated
// in a provided build root (represented by another pipeline). All inputs
// to a pipeline must be explicitly specified, either in terms of other
// pipeline, in terms of content addressable inputs or in terms of static
// parameters to the inherited Pipeline structs.
package manifest

import (
	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/container"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/ostree"
	"github.com/osbuild/osbuild-composer/internal/platform"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

type Pipeline interface {
	Name() string
	Checkpoint()
	Export() *artifact.Artifact
	getCheckpoint() bool
	getExport() bool
	getBuildPackages() []string
	getPackageSetChain() []rpmmd.PackageSet
	serializeStart([]rpmmd.PackageSpec)
	serializeEnd()
	serialize() osbuild.Pipeline
	getPackageSpecs() []rpmmd.PackageSpec
	getOSTreeCommits() []ostree.CommitSpec
	getContainerSpecs() []container.Spec
	getInline() []string
}

// A Base represents the core functionality shared between each of the pipeline
// implementations, and the Base struct must be embedded in each of them.
type Base struct {
	manifest   *Manifest
	name       string
	build      *Build
	checkpoint bool
	export     bool
}

// Name returns the name of the pipeline. The name must be unique for a given manifest.
// Pipeline names are used to refer to pipelines either as dependencies between pipelines
// or for exporting them.
func (p Base) Name() string {
	return p.name
}

func (p *Base) Checkpoint() {
	p.checkpoint = true
}

func (p Base) getCheckpoint() bool {
	return p.checkpoint
}

func (p *Base) Export() *artifact.Artifact {
	panic("can't export pipeline")
}

func (p Base) getExport() bool {
	return p.export
}

func (p Base) GetManifest() *Manifest {
	return p.manifest
}

func (p Base) getBuildPackages() []string {
	return []string{}
}

func (p Base) getPackageSetChain() []rpmmd.PackageSet {
	return nil
}

func (p Base) getPackageSpecs() []rpmmd.PackageSpec {
	return []rpmmd.PackageSpec{}
}

func (p Base) getOSTreeCommits() []ostree.CommitSpec {
	return nil
}

func (p Base) getContainerSpecs() []container.Spec {
	return nil
}

func (p Base) getInline() []string {
	return []string{}
}

// NewBase returns a generic Pipeline object. The name is mandatory, immutable and must
// be unique among all the pipelines used in a manifest, which is currently not enforced.
// The build argument is a pipeline representing a build root in which the rest of the
// pipeline is built. In order to ensure reproducibility a build pipeline must always be
// provided, except for int he build pipeline itself. When a build pipeline is not provided
// the build host's filesystem is used as the build root. The runner specifies how to use this
// pipeline as a build pipeline, by naming the distro it contains. When the host system is used
// as a build root, then the necessary runner is autodetected.
func NewBase(m *Manifest, name string, build *Build) Base {
	p := Base{
		manifest: m,
		name:     name,
		build:    build,
	}
	if build != nil {
		if build.Base.manifest != m {
			panic("build pipeline from a different manifest")
		}
	}
	return p
}

// serializeStart must be called exactly once before each call
// to serialize().
func (p Base) serializeStart([]rpmmd.PackageSpec) {
}

// serializeEnd must be called exactly once after each call to
// serialize().
func (p Base) serializeEnd() {
}

// Serialize turns a given pipeline into an osbuild.Pipeline object. This object is
// meant to be treated as opaque and not to be modified further outside of the pipeline
// package.
func (p Base) serialize() osbuild.Pipeline {
	pipeline := osbuild.Pipeline{
		Name: p.name,
	}
	if p.build != nil {
		pipeline.Build = "name:" + p.build.Name()
	}
	return pipeline
}

type Tree interface {
	Name() string
	GetManifest() *Manifest
	GetPlatform() platform.Platform
}
