package semver

import (
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/pkg/errors"
)

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "parsing tag"
}

// genReleaseTag prints release tag to the console for easy git tagging.
func releaseTag(ctx *context.Context, version string) string {
	relPrefix := "RELEASE"
	if ctx.Snapshot {
		relPrefix = "DEVELOPMENT"
	}

	relTag := strings.Replace(version, " ", "-", -1)
	relTag = strings.Replace(relTag, ":", "-", -1)
	relTag = strings.Replace(relTag, ",", "", -1)
	return relPrefix + "." + relTag
}

// Run executes the hooks.
func (Pipe) Run(ctx *context.Context) error {
	if ctx.GenerateMinIO {
		buildTime := time.Now().UTC().Format(time.RFC3339)
		ctx.MinIO = context.MinIOInfo{
			Version:    buildTime,
			ReleaseTag: releaseTag(ctx, buildTime),
		}
		ctx.Git.CurrentTag = ctx.MinIO.ReleaseTag
		return nil
	}

	sv, err := semver.NewVersion(ctx.Git.CurrentTag)
	if err != nil {
		if ctx.Snapshot {
			return pipe.ErrSnapshotEnabled
		}
		if ctx.SkipValidate {
			log.WithError(err).
				WithField("tag", ctx.Git.CurrentTag).
				Warn("current tag is not a semantic tag")
			return pipe.ErrSkipValidateEnabled
		}
		return errors.Wrapf(err, "failed to parse tag %s as semver", ctx.Git.CurrentTag)
	}
	ctx.Semver = context.Semver{
		Major:      sv.Major(),
		Minor:      sv.Minor(),
		Patch:      sv.Patch(),
		Prerelease: sv.Prerelease(),
	}
	return nil
}
