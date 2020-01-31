package version

// These should be overridden at build-time with
// `LDFLAGS=-ldflags "-X '$(go list)/version.Version=$VERSION' -X '$(go list)/version.BuildTime=$BUILDTIME'"`.
// The suggested values of `$VERSION` and `$BUILDTIME` are `git describe --tags` and `date -u +"%Y-%m-%dT%H:%M:%SZ"`,
// respectively.
var Version = "dev"
var BuildTime = "unknown"
