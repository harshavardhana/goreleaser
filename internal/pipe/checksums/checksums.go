// Package checksums provides a Pipe that creates .checksums files for
// each artifact.
package checksums

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Pipe for checksums.
type Pipe struct{}

func (Pipe) String() string {
	return "calculating checksums"
}

// Default sets the pipe defaults.
func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Checksum.NameTemplate == "" {
		ctx.Config.Checksum.NameTemplate = "{{ .ProjectName }}_{{ .Version }}_checksums.txt"
	}
	if ctx.Config.Checksum.Algorithm == "" {
		ctx.Config.Checksum.Algorithm = "sha256"
	}
	return nil
}

// Run the pipe.
func (Pipe) Run(ctx *context.Context) (err error) {
	artifactList := ctx.Artifacts.Filter(
		artifact.Or(
			artifact.ByType(artifact.UploadableArchive),
			artifact.ByType(artifact.UploadableBinary),
			artifact.ByType(artifact.UploadableSourceArchive),
			artifact.ByType(artifact.LinuxPackage),
		),
	).List()
	if len(artifactList) == 0 {
		return nil
	}

	for _, arf := range artifactList {
		if err = checksums(ctx.Config.Checksum.Algorithm, arf); err != nil {
			return err
		}
		ctx.Artifacts.Add(&artifact.Artifact{
			Type: arf.Type,
			Path: arf.Path + "." + ctx.Config.Checksum.Algorithm + "sum",
			Name: arf.Name + "." + ctx.Config.Checksum.Algorithm + "sum",
		})
	}
	return nil
}

func checksums(algorithm string, artifact *artifact.Artifact) error {
	log.WithField("file", artifact.Name).Info("checksumming")
	sha, err := artifact.Checksum(algorithm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(artifact.Path+"."+algorithm+"sum", []byte(fmt.Sprintf("%v  %v\n", sha, filepath.Base(artifact.Name))), 0644)
}
